# chatAPI
Gorilla Websocket을 발전시켜서 API로 만들어 봅시다.

[라인 메세징 시스템 참고](https://engineering.linecorp.com/ko/blog/the-architecture-behind-chatting-on-line-live/)

## 실행
- 프로젝트 Root directory에 다음과 같이 `config.json` 파일 생성
   ```json
   {
    "db": {
      "id": "",
      "password": "",
      "host": "",
      "port": "",
      "database": "",
      "listeningChannelName": ""
    },
    "mq": {
      "host": "",
      "port": "",
      "id": "",
      "password": "",
      "listeningQueueName": ""
    },
    "redis": {
      "id":"",
      "password": "",
      "host": "",
      "port": "",
      "database": ""
    },
    "logger": {
      "level": ""
    }
  }
   ```
