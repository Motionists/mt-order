#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
将拍扁的文件名（例如 web_src_pages_Menu_Version2.vue）恢复到目录结构
目标示例：
  - server_internal_handlers_auth_Version2.go  -> server/internal/handlers/auth.go
  - server_migrations_001_init_Version2.sql    -> server/migrations/001_init.sql
  - web_src_pages_Menu_Version2.vue            -> web/src/pages/Menu.vue
  - web_Dockerfile_Version2                    -> web/Dockerfile
  - web_vite.config_Version2.ts                -> web/vite.config.ts
"""
import argparse
import os
import re
import shutil
import subprocess
from typing import List, Tuple

def list_flat_files(cwd: str) -> List[str]:
    files = []
    for fn in os.listdir(cwd):
        # 只处理当前目录下的文件（避免已在子目录中的正常文件）
        if not os.path.isfile(fn):
            continue
        # 带 _Version2 或 _Version2.<ext> 的文件
        if re.search(r"_Version2(\.[^.]+)?$", fn):
            files.append(fn)
    return files

def build_target_path(src: str) -> str:
    """
    规则：
    - 去掉最后一个 "_Version2"
    - 首段必须是 "server" 或 "web"
    - 其余用下划线拆分成目录，最后一段作为文件名（附回原扩展名）
    - 特判 server_migrations_001_xxx.sql -> server/migrations/001_xxx.sql
    """
    base, ext = os.path.splitext(src)
    # 去掉末尾一次 _Version2（不影响 ext）
    base = re.sub(r"_Version2$", "", base)

    parts = base.split("_")
    if parts[0] not in ("server", "web"):
        # 非目标前缀，保持原样到 misc 目录，避免丢失
        return os.path.join("misc", src)

    root = parts[0]

    # 特判 migrations: server_migrations_001_init.sql
    if root == "server" and len(parts) >= 4 and parts[1] == "migrations":
        # 目录部分：server / migrations / (可能还有子目录)
        # 文件名：倒数第二段 + "_" + 倒数第一段 + 扩展名
        # 例如: ['server','migrations','001','init'] + '.sql' -> 001_init.sql
        dir_parts = parts[1:-2]  # 包含 'migrations'
        filename = f"{parts[-2]}_{parts[-1]}{ext}"
        return os.path.join(root, *dir_parts, filename)

    # 普通规则：最后一段为文件名，其余为目录
    dir_parts = parts[1:-1]
    filename = f"{parts[-1]}{ext}"
    return os.path.join(root, *dir_parts, filename)

def ensure_dir(path: str):
    d = os.path.dirname(path)
    if d and not os.path.exists(d):
        os.makedirs(d, exist_ok=True)

def git_mv(src: str, dst: str) -> bool:
    # 在 git 仓库中优先使用 git mv，保留历史
    if not os.path.isdir(".git"):
        return False
    try:
        subprocess.run(["git", "mv", "-f", src, dst], check=True, capture_output=True)
        return True
    except subprocess.CalledProcessError:
        return False

def move_file(src: str, dst: str, force: bool):
    ensure_dir(dst)
    if os.path.exists(dst):
        if force:
            # 若目标存在，先删除（git rm 或直接删）
            if os.path.isdir(".git"):
                try:
                    subprocess.run(["git", "rm", "-f", dst], check=True, capture_output=True)
                except subprocess.CalledProcessError:
                    os.remove(dst)
            else:
                os.remove(dst)
        else:
            print(f"[skip] target exists: {dst}")
            return

    if not git_mv(src, dst):
        shutil.move(src, dst)
    print(f"[moved] {src} -> {dst}")

def main():
    ap = argparse.ArgumentParser(description="Reshape flattened files into proper directories.")
    ap.add_argument("--apply", action="store_true", help="Actually move files (default is dry run).")
    ap.add_argument("--force", action="store_true", help="Overwrite existing targets.")
    args = ap.parse_args()

    files = list_flat_files(".")
    if not files:
        print("No *_Version2 files found in current directory.")
        return

    plan: List[Tuple[str, str]] = []
    for f in sorted(files):
        dst = build_target_path(f)
        plan.append((f, dst))

    print("Planned moves:")
    for src, dst in plan:
        print(f"  {src}  ->  {dst}")

    if not args.apply:
        print("\nDry run only. Add --apply to perform moves.")
        return

    for src, dst in plan:
        move_file(src, dst, args.force)

    print("\nDone. Review changes, then:")
    if os.path.isdir(".git"):
        print("  git status")
        print('  git add -A && git commit -m "chore: reshape files into tree and drop _Version2 suffix"')
    else:
        print("  (not a git repo)")

if __name__ == "__main__":
    main()