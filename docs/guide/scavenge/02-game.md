---
order: 2
---

# The Scavenger Hunt Game

The application you are building today can be used in many different ways but I'll be talking about building the app as a **scavenger hunt** game. Scavenger hunts are all about someone setting up tasks or questions that challenge a participant to find solutions that come with some sort of a prize. The basic mechanics of the game are as follows:

* Anyone can post a question with an encrypted answer.
* This question comes paired with a bounty of coins.
* Anyone can post an answer to this question. If the answer is correct, that person receives the bounty of coins.

Something to note here is that when dealing with a public network with latency, it is possible that something like a [man-in-the-middle attack](https://en.wikipedia.org/wiki/Man-in-the-middle_attack) could take place. Instead of pretending to be one of the parties, an attacker would take the sensitive information from one party and use it for their own benefit. This scenario is actually called [Front Running](https://en.wikipedia.org/wiki/Front_running) and happens as follows:

1. You post the answer to some question with a bounty attached to it.
2. Someone else sees you posting the answer and posts it themselves right before you.
3. Since they posted the answer first, they receive the reward instead of you.

To prevent Front-Running, you will implement a **commit-reveal** scheme. A commit-reveal scheme converts a single exploitable interaction and turns it into two safe interactions.

**The first interaction is the commit**. This is where you "commit" to posting an answer in a follow-up interaction. This commit consists of a cryptographic hash of your name combined with the answer that you think is correct. The app saves that value which is a claim that you know the answer but that it hasn't been confirmed whether the answer is correct.

**The next interaction is the reveal**. This is where you post the answer in plaintext along with your name. The application will take your answer and your name and cryptographically hash them. If the result matches what you previously submitted during the commit stage, then it will be proof that it is in fact you who knows the answer, and not someone who is just front-running you.

A system like this could be used in tandem with any kind of gaming platform in a **trustless** way. Imagine you were playing "The Legend of Zelda" and the game was compiled with all the answers to different scavenger hunts already included. When you beat a level the game could reveal the secret answer. Then either explicitly or behind the scenes, this answer could be combined with your name, hashed, submitted and subsequently revealed. Your name would be rewarded and you would have more points in the game.

Another way of achieving this level of security would be to have an Access Control List where there was an admin account that the video game company controlled. This admin account could confirm that you beat the level and then give you points. The problem with this is that it creates a ***single point of failure*** and a single target for trying to attack the system. If there is one key that rules the castle then the whole system is broken if that key is compromised. It also creates a problem with coordination if that Admin account has to be online all the time in order for players to get their points. If you use a commit reveal system then you have a more trustless architecture where you don't need permission to play. This design decision has benefits and drawbacks, but paired with a careful implementation it can allow your game to scale without a single bottle neck or point of failure.

Now that we know what we're building we can get started.
