---
name: open-pr
description: Abre um Pull Request no GitHub respeitando o template do projeto gymflow-api. Informe o número da task Jira (ex: GYM-50) e uma descrição do que foi feito.
---

Você é responsável por abrir Pull Requests no repositório gymflow-api seguindo exatamente o template em `.github/pull_request_template.md`.

## O que fazer

1. Leia o template em `.github/pull_request_template.md`
2. Rode `git diff main...HEAD` e `git log main..HEAD --oneline` para entender o que mudou
3. Preencha o template com base nas mudanças reais do diff:
   - **O que esse PR faz?** — descreva o que foi implementado em 2-3 frases
   - **Task relacionada** — substitua `GYM-XXX` pelo número real da task nos três lugares (link Jira, link Confluence e `Closes`). O link Jira segue o padrão `https://gabrielamiranda1110.atlassian.net/browse/GYM-XXX`. Para o Confluence, busque a página da task usando `mcp__claude_ai_Atlassian__searchConfluenceUsingCql` com `text = "GYM-XXX"` — se não encontrar, omita o link do Confluence.
   - **Como testar** — liste os passos para testar manualmente (seja específico: comandos reais, endpoints, etc.)
   - **Checklist** — marque apenas os itens que realmente se aplicam à mudança

4. Abra o PR com `gh pr create` usando o template preenchido
5. Aplique a label correta com `gh pr edit <número> --add-label <label>`. Labels disponíveis: `feat`, `fix`, `chore`, `refactor`, `documentation`
6. Retorne a URL do PR

## Regras

- Nunca invente passos de teste — baseie-se no que realmente foi implementado
- Não marque itens do checklist que não se aplicam (ex: "Migration criada" se não houve mudança no banco)
- Não remova seções do template — preencha todas
- Use `--base main` no `gh pr create`
- O título do PR deve seguir conventional commits: `tipo: descrição curta` (ex: `feat: criar endpoint de treino`)
