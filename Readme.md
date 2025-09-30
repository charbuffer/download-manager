# Download Manager

This project is a REST API service which accepts URLs in incoming Tasks and downloads files from them into **Downloads** folder
Each Task created by request contains list of files with statuses

### To launch application use: 

`make run`

## Used patterns and concepts

- Dependency Inversion
- Worker Pool
- Graceful Shutdown
- Config
- Layered Architecture
- Repository Pattern
- DTO

## Routes

`GET /task` - Get all tasks

`GET /task/:id` - Get task by id

`POST /task` - Add new task

Example request body:
```json
{
    "urls": [
        "https://example.com/", 
        "https://example.com"
    ]
}
```