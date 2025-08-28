# Live Chat Demo
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

## Local Dev Related

### Install and Build Env
docker compose build --no-cache
docker compose up -d

### Dev Env Debugging
docker compose logs -f app

### Remove Env
docker compose down -v

---

## Todo List 📝

### Dev Environment
- [X] Local dev env

### User Features
- [ ] Member login function

### Cloud / Deployment
- [X] GCP cloud run set up
- [X] Mongo Atlas(GCP) set up
- [X] Build project on GCP cloud run

### AI Features
- [ ] AI feature research
- [ ] Add some AI feature to project

### Security
- [ ] XSS prevention (Content Security Policy, escaping)
- [ ] CSRF protection
- [ ] Authentication & authorization checks

### Tools / Maintenance
- [ ] Others (linter, CI/CD, etc.)


