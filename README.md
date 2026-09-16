<h1 align="center">
    Irterlan RP | Telegram Bot
</h1>
<p align="center">
    <img alt="Go version" src="https://img.shields.io/github/go-mod/go-version/likcoah/irterlan-rp-bot">
    <img alt="Docker & Compose" src="https://img.shields.io/badge/Dockerfile_%26_Compose-2496ED?logo=docker&logoColor=white">
    <a href="https://t.me/IrterlanRPBot"><img alt="Telegram Bot" src="https://img.shields.io/badge/Telegram_Bot-26A5E4?logo=telegram&logoColor=white"></a>
</p>

> [!NOTE]
> **Work in Progress (WIP)**
> Project under development

### About
**Irterlan RP Bot** is a custom bot for the Irterlan Telegram role-playing group, designed to use a custom dice-rolling system by likcoah and simplify gameplay

### Planned features
- [ ] **Universal Mathematical & Dice Parser**: A single engine for evaluating math expressions and rolling dice directly in chat.
  - Supports custom dice syntax like `/1d20 + 2d6 - 5` with support for Cyrillic `д` (`/1д20`).
  - Evaluates standard arithmetic expressions (`/5 + 10 * 2 - 90`).
  - Uses an AST (Abstract Syntax Tree) parser under the hood for clean operator precedence and secure execution.

### Features under consideration
- Character sheets and interactive cell selection within them
- Automate battles using combat mode
- AI integration? (xd)
