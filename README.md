# Language Learning Platform

Uma plataforma educacional SaaS para ensino de idiomas com arquitetura de microserviços.

## 🎯 Funcionalidades

- **Para Professores**: Criar cursos, tarefas, correção e feedback
- **Para Alunos**: Receber tarefas, estudar, acompanhar progresso
- **Consultoria Online**: Aulas ao vivo, acompanhamento individualizado
- **Multi-idioma e Multi-professor**: Arquitetura SaaS escalável

## 🏗️ Arquitetura

### Camadas
- **Frontend**: Next.js (web) + React Native (mobile)
- **API Gateway**: NGINX/Traefik com autenticação e rate limiting
- **Microserviços**: C# (Auth) + Go (demais serviços)
- **Infraestrutura**: PostgreSQL, Redis, S3/MinIO, Observabilidade

### Serviços
- `auth-service` (C#): Autenticação e autorização
- `user-service` (Go): Gestão de usuários
- `course-service` (Go): Gestão de cursos
- `task-service` (Go): Tarefas e exercícios
- `progress-service` (Go): Acompanhamento de progresso
- `notification-service` (Go): Notificações
- `file-service` (Go): Upload e gestão de arquivos
- `video-service` (Go): Integração de vídeo conferência

## 🚀 Quick Start

### Desenvolvimento Local
```bash
# Iniciar todos os serviços
./scripts/local-dev.sh

# Executar migrações
./scripts/migrate.sh

# Popular banco de dados
./scripts/seed-db.sh
```

### Docker Compose
```bash
cd infra/docker
docker-compose up -d
```

### Kubernetes
```bash
# Deploy em dev
kubectl apply -k infra/kubernetes/overlays/dev/

# Deploy em production
kubectl apply -k infra/kubernetes/overlays/production/
```

## 📁 Estrutura do Projeto

```
language-platform/
├── frontend/          # Aplicações frontend
├── services/          # Microserviços
├── libs/              # Bibliotecas compartilhadas
├── infra/             # Infraestrutura e deployment
├── scripts/           # Scripts auxiliares
└── docs/              # Documentação
```

## 🛠️ Tecnologias

- **Backend**: C# (.NET 8), Go 1.21+
- **Frontend**: Next.js 14, React Native
- **Database**: PostgreSQL 15+
- **Cache**: Redis 7+
- **Storage**: MinIO / S3
- **Observability**: Prometheus, Grafana, ELK, Jaeger
- **Container**: Docker, Kubernetes

## 📚 Documentação

- [Arquitetura](docs/architecture.md)
- [Contratos de API](docs/api-contracts.md)
- [Decisões Técnicas](docs/decisions.md)

## 🗂️ Organização de Arquivos

- Documentacao adicional esta em `docs/`
- Arquivos de configuracao e solucao estao em `config/`

## 📝 Licença

MIT License - veja LICENSE para detalhes
