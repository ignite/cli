---
sidebar_position: 2
---

# Scavenger hunt game

This tutorial focuses on building the app as a **scavenger hunt** game. Scavenger hunts are all about someone setting up tasks or questions that challenge a participant to find solutions that come with a prize. The basic mechanics of the game are as follows:

* Anyone can post a question with an encrypted answer.
* This question comes paired with a bounty of coins.
* Anyone can post an answer to this question. If the answer is correct, that person receives the bounty of coins.

## Safe interactions

On a public network with latency, it is possible that something like a [man-in-the-middle attack](https://en.wikipedia.org/wiki/Man-in-the-middle_attack) could take place. Instead of pretending to be one of the parties, an attacker takes the sensitive information from one party and uses it for their own benefit. This scenario is called [Front Running](https://en.wikipedia.org/wiki/Front_running) and happens as follows:

1. You post the answer to some question with a bounty attached to it.
2. Someone else sees you posting the answer and posts it themselves right before you.
3. Since they posted the answer first, they receive the reward instead of you.

### Prevent front running

To prevent front running, implement a commit-reveal scheme that converts a single exploitable interaction into two safe interactions.

The first interaction is the commit where you "commit" to posting an answer in a follow-up interaction. This commit consists of a cryptographic hash of your name combined with the answer that you think is correct. The app saves that value as a claim that you know the answer, but hasn't yet confirmed whether the answer is correct.

The second interaction is the reveal where you post the answer in plain text along with your name. The application takes your answer and your name and cryptographically hashes them. If the result matches what you previously submitted during the commit stage, that is the proof that it is in fact you who knows the answer and not someone who is just front-running you.

### Security

A system like this commit-reveal scheme could be used in tandem with any kind of gaming platform in a **trustless** way. Imagine playing "The Legend of Zelda" and the game was compiled with all the answers to different scavenger hunts already included. When you obtain a certain level the game could reveal the secret answer. Then explicitly or behind the scenes, this answer could be combined with your name, hashed, submitted, and subsequently revealed. Your name could be rewarded so that you gain more points in the game.

Another way of achieving this level of security is with access-control list (ACL) managed by an admin account under control of the gaming company. This admin account could confirm that you beat the level and then give you points. The problem with this ACL approach is a ***single point of failure*** and a single target for trying to attack the system. If there is one key that rules the castle then the whole system is broken if that key is compromised. ACL security also creates a problem with coordination if that admin account is required to be online in order for players to get their points. With a commit-reveal system you have a more trustless architecture where permission is not required to play. The commit-reveal design decision has benefits and drawbacks, but when paired with a careful implementation your game can scale without a single bottleneck or point of failure.

Now that you know what you are building, you can get started building your game.
