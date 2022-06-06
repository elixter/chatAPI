# chatAPI
Gorilla Websocket을 발전시켜서 API로 만들어 봅시다.

[라인 메세징 시스템 참고](https://engineering.linecorp.com/ko/blog/the-architecture-behind-chatting-on-line-live/)

## 기술스택
- Go
- MySQL
- Redis (pub/sub)

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
      "database": "",
      "listeningChannelName": "",
      "publishChannelName": ""
    },
    "logger": {
      "level": ""
    }
  }
   ```

## Demo
![Hnet-image](https://user-images.githubusercontent.com/45375508/172126121-641e86e8-d674-46f4-b39c-3ecc7da213df.gif)
