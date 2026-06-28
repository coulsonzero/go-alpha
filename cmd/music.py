#!/usr/bin/env python3
"""
探测腾讯云COS music目录下的所有MP3文件并输出JSON列表
"""

import json
import sys
from urllib.parse import quote

import requests

BASE_URL = "https://music-1309877840.cos.ap-beijing.myqcloud.com/music/"


def list_mp3_files() -> list[dict]:
    """依次尝试多种方式获取歌曲列表"""
    # 尝试各种索引文件
    index_urls = [
        f"{BASE_URL}index.json",
        f"{BASE_URL}list.json",
        f"{BASE_URL}songs.json",
        f"{BASE_URL}music.json",
    ]
    for url in index_urls:
        try:
            resp = requests.get(url, timeout=10)
            if resp.status_code == 200:
                data = resp.json()
                items = data if isinstance(data, list) else \
                        data.get("data") if isinstance(data, dict) else None
                if items and isinstance(items, list):
                    return parse_songs(items)
        except Exception:
            continue
    return []


def parse_songs(items: list) -> list[dict]:
    results = []
    for item in items:
        if isinstance(item, str):
            name = item.replace(".mp3", "").strip()
            parts = name.split(" - ", 1)
            if len(parts) == 2:
                artist, title = parts
            else:
                title = name
                artist = "未知"
            results.append({
                "title": title.strip(),
                "artist": artist.strip(),
                "src": f"{BASE_URL}{quote(item)}",
            })
        elif isinstance(item, dict):
            title = item.get("title", item.get("name", "未知"))
            artist = item.get("artist", item.get("author", "未知"))
            src = item.get("src", item.get("url", ""))
            if not src:
                filename = item.get("file", item.get("filename", f"{title}.mp3"))
                src = f"{BASE_URL}{quote(filename)}"
            results.append({
                "title": title.strip(),
                "artist": artist.strip(),
                "src": src,
            })
    return results


def main():
    try:
        results = list_mp3_files()
        print(json.dumps({"code": 200, "data": results}, ensure_ascii=False))
    except Exception as e:
        print(json.dumps({"code": 500, "error": str(e)}, ensure_ascii=False))


if __name__ == "__main__":
    main()
