#  LLM with Go — MCP Hub Experimental

> Hub modular em Go para conectar **LLMs**, **ferramentas de dev** e **fontes de dados** via **Model Context Protocol (MCP)**.

---

## Visão

Projeto em construção que busca criar um **MCP Hub** é um ponto central onde llms e ferramentas  podem interagir de forma integrada e local.

O objetivo é permitir que desenvolvedores executem tarefas reais (consultas, automações, resumos, análises) através de um protocolo aberto e extensível.

---

## Estado Atual

- Protótipo funcional em Go com **comunicação cliente-servidor local** via `mcp.NewInMemoryTransports()`.
- Ferramenta `llm-caller` que recebe uma pergunta e retorna resposta de um wrapper LLM.
- Logging básico e estrutura modular pronta para expansão.

---

## Próximos Passos

- Suporte a transportes reais (**Unix/TCP/gRPC**).
- Integração com provedores de LLM (OpenAI, Ollama, etc).
- Novas ferramentas MCP.
- Configuração externa via `.env`.
- CLI e API Gateway.
- Logs estruturados e métricas Prometheus.