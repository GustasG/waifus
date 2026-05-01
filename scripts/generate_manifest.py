import argparse
import json
import os
import sys
from pathlib import Path

IMAGE_EXTENSIONS = {".jpg", ".jpeg", ".png", ".gif", ".webp", ".avif"}


def generate(languages_dir: Path) -> dict[str, list[str]]:
    manifest = {}

    for lang in sorted(os.listdir(languages_dir)):
        lang_path = languages_dir / lang

        if not lang_path.is_dir():
            continue

        images = sorted(
            f
            for f in os.listdir(lang_path)
            if os.path.splitext(f)[1].lower() in IMAGE_EXTENSIONS
        )

        if images:
            manifest[lang] = images

    return manifest


def main():
    parser = argparse.ArgumentParser(
        description="Generate image manifest from a languages directory.",
    )

    parser.add_argument(
        "-d",
        "--directory",
        default=Path("assets", "languages"),
        help="Path to the languages directory",
        type=Path,
    )

    parser.add_argument(
        "-o",
        "--output",
        default=Path("manifest.json"),
        help="Output JSON file path",
        type=Path,
    )

    args = parser.parse_args()

    if not args.directory.is_dir():
        print(f"error: '{args.directory}' is not a directory", file=sys.stderr)
        sys.exit(1)

    manifest = generate(args.directory)

    with open(args.output, "w", encoding="utf-8") as f:
        json.dump(manifest, f, ensure_ascii=False)

    print(
        f"wrote {args.output}: {len(manifest)} languages, {sum(len(v) for v in manifest.values())} images"
    )


if __name__ == "__main__":
    main()
