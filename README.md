<p align="center"><img width="658" alt="eXit Intro" src="https://github.com/user-attachments/assets/75800f4d-dbc1-4066-8641-d6875840cddb"></p>

# eXit - A Terminal-based Game Recreation from *Mr. Robot*

Welcome to **eXit**, a faithful recreation of the terminal-based game featured in *Mr. Robot* (S4E11). This project is built using [Go](https://golang.org) and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework for a nostalgic yet modern terminal experience. 

Not only does this game stay true to the original look and feel from the show, but it also includes hidden **easter eggs** that fans of *Mr. Robot* will appreciate. Additionally, I've utilized [Ansify](https://github.com/ZachLTech/Ansify), one of my other projects, to convert key scenes from *Mr. Robot* into terminal-friendly ANSI images, creating an immersive and authentic gameplay environment.

Soon, the game will support SSH access via Docker containers, allowing anyone to play directly from their own terminal!

## 🚀 Features

- **Classic Terminal-based Gameplay**: Experience the eerie and nostalgic feeling of playing *eXit* just like in the show.
- **Hidden Easter Eggs**: Discover references to iconic scenes and moments from *Mr. Robot* as you play.
- **ANSI Art Integration**: Utilizes [Ansify](https://github.com/ZachLTech/Ansify) to bring scenes from *Mr. Robot* to life in the terminal with ANSI images.
- **Terminal Animations**: Smooth terminal animations enhance the gaming experience.
- **Future Terminal Game Engine**: This project will serve as a foundation for a fully-fledged terminal-based game engine in the future.
- **SSH Support Coming Soon**: Play the game over SSH with Docker integration for a seamless, accessible experience.

## 📦 Installation

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.19+)
- [Git](https://git-scm.com/)
- (Optional) [Docker](https://www.docker.com/get-started) - for the upcoming SSH version

### Clone the Repository

```bash
git clone https://github.com/ZachLTech/eXit.git
cd eXit
```

### Build the Game

```bash
go build -o exit
```

### Run the Game

```bash
./exit
```

## 🛠 Development

### Dependencies

To install dependencies:

```bash
go mod tidy
```

### Run in Development Mode

```bash
go run main.go
```

## 🎨 ANSI Art & Animations

This project integrates with [Ansify](https://github.com/ZachLTech/Ansify) to convert scenes from *Mr. Robot* into terminal-friendly ANSI images. As you progress through the game, you’ll encounter animations and visuals that bring the *Mr. Robot* universe directly into your terminal.

### Creating Your Own ANSI Art

If you'd like to customize the game’s scenes or add your own, you can use [Ansify](https://github.com/ZachLTech/Ansify) to convert images into ANSI format, ensuring they are compatible with the terminal interface.

## 🖥️ SSH Access (Coming Soon!)

I’m currently working on Docker-based SSH support so you can play *eXit* directly from your terminal.
Stay tuned for updates on SSH support!

## 🌌 Future Plans

- **Terminal Game Engine**: I'm working towards developing a full-featured terminal-based game engine based on the mechanics and aesthetics of *eXit*. This engine will support animations, ANSI art, and complex game mechanics—all within a terminal environment.
- **More Easter Eggs**: Expect additional hidden features and references for *Mr. Robot* enthusiasts... and maybe this will be a key to the puzzle for my CTF/Treasure Hunt (Coming Soon...)

## 🤝 Contributing

Contributions are welcome! Feel free to fork this project, make your improvements, and submit a pull request.

## 💬 Feedback & Support

If you encounter any issues, have questions, or want to share your thoughts, feel free to open an [issue](https://github.com/ZachLTech/eXit/issues) or reach out to me directly!

## ⭐ Acknowledgments

- Inspired by *Mr. Robot* and its iconic terminal-based game.
- Built using the fantastic [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.
- Thanks to the open-source community for making projects like this possible.
- Special shoutout to my other project [Ansify](https://github.com/ZachLTech/Ansify) for making ANSI art integration possible.