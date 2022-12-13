# How To [Build/Create/Do Something] in Ignite CLI

<!--
Use this tutorial template as a quick starting point when writing Ignite CLI how-to tutorials. 

After you review the template, delete the comments and begin writing your outline or article. Examples of Markdown formatting syntax are provided at the bottom of this template.

As you write, refer to industry style and formatting guidelines. 

We admire, respect, and rely on these resources:

- Google developer documentation [style guide](https://developers.google.com/style)
- Digital Ocean style guide [do.co/style](https://do.co/style)

[Create an issue](https://github.com/ignite/cli/issues/new/choose) to let us know if you have questions. 

-->

<!-- To preview a content deploy so you can see what your article looks like before it is published, see [../CONTRIBUTING.md#viewing-tutorial-builds]. 

Our users must be able to follow the tutorial from beginning to end on their own computer. Before submitting a tutorial for PR review, be sure to test the content by completing all steps from start to finish exactly as they are written. Cut and paste commands from the article into your terminal to make sure that typos are not present in the commands. If you find yourself executing a command that isn't in the article, incorporate that command into the article to make sure the user gets the exact same results. 
-->

<!-- Use sentence case for all headings and titles, see https://capitalizemytitle.com/ -->

<!-- Use GitHub flavored Markdown, see [Mastering Markdown](https://docs.github.com/en/get-started/writing-on-github/getting-started-with-writing-and-formatting-on-github/basic-writing-and-formatting-syntax)  -->

<!-- Our articles have a specific structure. How-To tutorials follow this structure:

* Front matter metadata
* Title
* Introduction and purpose (Level 2 heading)
* Prerequisites (Level 2 heading)
* Step 1 — Doing something (Level 2 heading)
* Step 2 — Doing something (Level 2 heading)
...
* Step 5 — Doing something (Level 2 heading)
* Conclusion (Level 2 heading)

 -->

### Introduction and purpose

Introductory paragraph about the topic that explains what this topic is about and why the user should care; what problem does the tutorial solve?

In this guide, you will [accomplish/build/] [some important thing]...

When you're finished, you'll be able to...

**Note:** The code in this tutorial is written specifically for this learning experience and is intended only for educational purposes. This tutorial code is not intended to be used in production.

## Prerequisites

<!-- Prerequisites let you leverage existing tutorials so you don't have to repeat installation or setup steps in your tutorial. -->

To complete this tutorial, you will need:

* A local development environment for [your chain] 
* Familiarity with the Cosmos ecosystem and [your chain]. See [cosmos.network](EIP-1559 for $ATOM) to learn more.
* (Optional) If software such as Git, Go, Docker, or other tooling needs to be installed, link to the proper article describing how to install it.
* (Optional) List any other accounts needed.

<!-- Example - uncomment to use

- A supported version of [Ignite CLI](/). To install Ignite CLI, see [Install Ignite CLI](../../guide/01-install.md). 
* A text editor like [Visual Studio Code](https://code.visualstudio.com/download).
* A web browser like [Chrome](https://www.google.com/chrome) or [Firefox](https://www.mozilla.org/en-US/firefox/new).

-->

## Step 1 — Doing something

Introduction to the step. What are you going to do and why are you doing it?

First....

Next...

Finally...

<!-- When showing a command, explain the command first by talking about what it does. Then show the command. Then show its output in a separate output block: -->

To verify the version of Ignite CLI that is installed, run the following command:

```bash
ignite --version
```

You'll see release details like the following output:

```
Ignite version:	v0.19.6
Ignite CLI build date:	2021-12-18T05:56:36Z
Ignite CLI source hash:	-
Your OS:		darwin
Your arch:		amd64
Your go version:	go version go1.16.4 darwin/amd64
```

<!-- When asking the user to open a file, be sure to specify the file name:

Create the `post.proto` file in your editor.

When showing the contents of a file, try to show only the relevant parts and explain what needs to change. -->

Modify the title by changing the contents of the `<title>` tag:

```protobuf
// ...

message Post {
  string creator = 1;
  string id = 2;
  string title = 3; 
  string body = 4; 
}

message MsgCreatePost {
  string creator = 1;
  string title = 2; 
  string body = 3; 
}

// ...
```

Now transition to the next step by telling the user what's next.

## Step 2 — Sentence case heading

Another introduction

Your content that guides the user to accomplish a specific step

Transition to the next step.

## Step 3 — Sentence case

Another introduction

Your content

Transition to the next step.

## Conclusion

In this article you [accomplished or built] [some important thing]. Now you can....

<!-- Speak to the benefits of this technique or procedure and optionally provide places for further exploration. -->

<!------------ Formatting ------------------------->

<!-- Some examples of how to mark up various things

This is _italics_ and this is **bold**.

Use italics and bold for specific things. 

This is `inline code`. Use single tick marks for filenames and commands.

Here's a command you can type on a command line:

```bash
which go
```

Here's output from a command:

```
/usr/local/go/bin/go
```

Write key presses in ALLCAPS.

Use a plus symbol (+) if keys need to be pressed simultaneously: `CTRL+C`.

**Note:** This is a note.

**Tip:** This is a tip.

Add diagrams and screenshots in PNG format with a self-describing filename. Embed them in the article using the following format:

![Alt text for screen readers](/path/to/img.png)

-->
