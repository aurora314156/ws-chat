# Simple Chat Demo
![Demo](assets/demo.png)

---

## Description

This is a simple real-time chat application deployed on **GCP**.  
It is built with **FastAPI**, **WebSocket**, and **MongoDB**, and allows all participants in the chat to see messages instantly.  

### Deployment & Services
- **Compute**: GCP Cloud Run  
- **Database**: MongoDB Atlas (hosted on GCP)  
- **Backend Framework**: FastAPI  
- **Real-time Communication**: WebSocket

---

## Dev Related

### Install and Build Env
docker compose up --build -d

### Dev Env Debugging
docker compose logs -f app

### Remove Env
docker compose down -v

---

## Todo List 📝

- [x] Local dev env
- [x] GCP cloud run set up
- [ ] Mongo Atlas(GCP) set up
- [ ] Build project on GCP cloud run
- [ ] AI feature research
- [ ] Add some AI feature to project
- [ ] Others (linter, CI/CD, etc.)