# README
## World Books
```base
views:
  - type: table
    name: 表格
    filters:
      and:
        - file.folder == "World"
    order:
      - file.name
      - tags

```

## Characters

```base
views:
  - type: table
    name: 表格
    filters:
      and:
        - file.folder == "Character"
    order:
      - file.name
      - name

```

## Chapters

```base
views:
  - type: table
    name: 表格
    filters:
      and:
        - file.folder == "Story"
    order:
      - file.name
      - title

```

