module.exports = {
  theme: "cosmos",
  title: "Starport",
  head: [
    [
      "script",
      {
        async: true,
        src: "https://www.googletagmanager.com/gtag/js?id=G-XL9GNV1KHW",
      },
    ],
    [
      "script",
      {},
      [
        "window.dataLayer = window.dataLayer || [];\nfunction gtag(){dataLayer.push(arguments);}\ngtag('js', new Date());\ngtag('config', 'G-XL9GNV1KHW');",
      ],
    ],
  ],
  themeConfig: {
    logo: {
      src: "/logo.png",
    },
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
              title: "Cosmos SDK Docs",
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
          "Chat with Starport and Cosmos SDK developers in <a href='https://discord.gg/H6wGTY8sxw' target='_blank'>Discord</a>.",
      },
      logo: "/logo.svg",
      textLink: {
        text: "starport.com",
        url: "https://starport.com/",
      },
      services: [
        {
          service: "medium",
          url: "https://medium.com/tendermint",
        },
        {
          service: "twitter",
          url: "https://twitter.com/starportHQ",
        },
        {
          service: "linkedin",
          url: "https://www.linkedin.com/company/tendermint/",
        },
        {
          service: "discord",
          url: "https://discord.gg/H6wGTY8sxw",
        },
        {
          service: "youtube",
          url: "https://www.youtube.com/channel/UCXMndYLK7OuvjvElSeSWJ1Q",
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
              url: "https://discord.gg/7fwqwc3afK",
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
