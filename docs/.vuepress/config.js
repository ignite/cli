module.exports = {
  theme: "cosmos",
  title: "Starport",
  themeConfig: {
    algolia: {
      id: "BH4D9OD16A",
      key: "d6908a9436133e03e9b0131bad808775",
      index: "docs-startport",
    },
    sidebar: {
      auto: true,
      nav: [
        {
          title: "Resources",
          children: [
            {
              title: "Starport on Github",
              path: "https://github.com/tendermint/starport",
            },
            {
              title: "Tutorials",
              path: "https://tutorials.cosmos.network",
            },
            {
              title: "Cosmos SDK docs",
              path: "https://docs.cosmos.network",
            },
          ],
        },
      ],
    },
    topbar: {
      banner: false,
    },
    custom: true,
    footer: {
      question: {
        text:
          "Chat with Starport and Cosmos SDK developers in <a href='https://discord.gg/W8trcGV' target='_blank'>Discord</a>.",
      },
      logo: "/logo.svg",
      textLink: {
        text: "cosmos.network/starport",
        url: "https://cosmos.network/starport",
      },
      services: [
        {
          service: "medium",
          url: "https://blog.cosmos.network/",
        },
        {
          service: "twitter",
          url: "https://twitter.com/cosmos",
        },
        {
          service: "linkedin",
          url: "https://www.linkedin.com/company/tendermint/",
        },
        {
          service: "reddit",
          url: "https://reddit.com/r/cosmosnetwork",
        },
        {
          service: "discord",
          url: "https://discord.gg/vcExX9T",
        },
        {
          service: "youtube",
          url: "https://www.youtube.com/c/CosmosProject",
        },
      ],

      smallprint:
        "This website is maintained by Tendermint Inc. The contents and opinions of this website are those of Tendermint Inc.",
      links: [
        {
          title: "Documentation",
          children: [
            {
              title: "Cosmos SDK",
              url: "https://docs.cosmos.network",
            },
            {
              title: "Cosmos Hub",
              url: "https://hub.cosmos.network",
            },
            {
              title: "Tendermint Core",
              url: "https://docs.tendermint.com",
            },
          ],
        },
        {
          title: "Community",
          children: [
            {
              title: "Cosmos blog",
              url: "https://blog.cosmos.network",
            },
            {
              title: "Forum",
              url: "https://forum.cosmos.network",
            },
            {
              title: "Chat",
              url: "https://discord.gg/W8trcGV",
            },
          ],
        },
        {
          title: "Contributing",
          children: [
            {
              title: "Contributing to the docs",
              url:
                "https://github.com/cosmos/cosmos-sdk/blob/master/docs/DOCS_README.md",
            },
            {
              title: "Source code on GitHub",
              url: "https://github.com/cosmos/cosmos-sdk/",
            },
          ],
        },
      ],
    },
  },
};
