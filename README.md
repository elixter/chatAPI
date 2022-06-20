# chatAPI
- [Gorilla Websocket의 채팅 예제](https://github.com/gorilla/websocket/tree/master/examples/chat)를 발전시켜서 API로 만들어 봅시다.

- [라인 메세징 시스템 참고](https://engineering.linecorp.com/ko/blog/the-architecture-behind-chatting-on-line-live/)

### 설명
라인에선 고속 병행 처리를 위해 Akka라는 것을 사용했으나 Go언어만을 사용해서 구축하였다. 각 구성요소에 해당하는 부분은 다음과 같다.
- at Line system => my system
   + UserActor => Client struct
   + ChatRoomActor => Room struct
   + ChatSupervisor => Hub struct

## 기술스택
- Go
- MySQL
- Redis

## TODO
- [x] 기본 기능
- [x] 서버간 동기화
- [x] 고루틴 누수 방지 코드 재설계
   + 널리 사용되는 패턴은 사용하지 않았지만 잔여 고루틴 없음
- [ ] 메세지 타입별 처리
   + [ ] 이미지 데이터 처리

- [ ] 채팅방 관리 API
- [ ] 커넥션 끊어졌을 경우 자동 재연결
   + [ ] 웹 소켓
   + [ ] Redis 
- [x] JSON 인코딩 벤치마크
   + 다른 방식으로 해봤지만 큰 차이가 없어서 기존 방식 사용

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
   
- 클라이언트에서 다음 URL로 소켓 연결
   ```cURL
   ws://hostname:8080/ws/:roomId
   ```
   + 예제 (javascript)
   ``` javascript
   // 연결
   conn = new WebSocket("ws://localhost:8080/ws/123");

   // 메세지 전송
   now = new Date(Date.now());
   conn.send(JSON.stringify({
       message_type: "chat_txt",
       author_id: 11,
       room_id: 123,
       content: msg.value,
       create_at: now.toISOString(),
   }));
   ```
   
## 소켓 데이터 포맷 (json)
```json
 {
   "message_type": message type : "chat_txt" or "chat_img"
   "author_id": integer value for user id
   "room_id": integer value for chatting room id
   "content": chatting content
   "create_at": UTC timestamp
 }
```


## Demo
![Hnet-image](https://user-images.githubusercontent.com/45375508/172126121-641e86e8-d674-46f4-b39c-3ecc7da213df.gif)
