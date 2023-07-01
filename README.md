# ArcSpy-Server

[ArcSpy](https://github.com/Lxns-Network/ArcSpy) 服务器端实现。

## Routes

### POST /player/sync

同步玩家数据。

#### JSON

| Key      | Type | Description |
|----------|------|-------------|
| `player` | dict | 玩家数据        |
| `scores` | list | 分数数据        |
| `cookie` | str  | 玩家 Cookie   |

### GET /player/data

获取玩家同步的玩家数据，需要 `Authorization` 头。

#### Parameters

| Name      | Type | Description |
|-----------|------|-------------|
| `user_id` | int  | 玩家 ID       |

### GET /player/scores

获取玩家同步的 Best 成绩列表，需要 `Authorization` 头。

#### Parameters

| Name      | Type | Description |
|-----------|------|-------------|
| `user_id` | int  | 玩家 ID       |

### GET /webapi/user/me

使用玩家同步的 Cookie 获取玩家数据，需要 `Authorization` 头。

#### Parameters

| Name      | Type | Description |
|-----------|------|-------------|
| `user_id` | int  | 玩家 ID       |

### GET /webapi/score/song/me

使用玩家同步的 Cookie 获取玩家单曲最佳成绩，需要 `Authorization` 头。

#### Parameters

| Name        | Type | Description |
|-------------|------|-------------|
| `user_id`   | int  | 玩家 ID       |
| `song_id`   | int  | 歌曲 ID       |
| `difficlty` | int  | 难度          |
