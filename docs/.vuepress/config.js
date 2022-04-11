module.exports = {
  theme: "cosmos",
  title: "Ignite CLI",
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
              title: "Ignite CLI on Github",
              path: "https://github.com/ignite-hq/cli",
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
          "Chat with Ignite CLI developers in <a href='https://discord.gg/ignt' target='_blank'>Discord</a>.",
      },
      logo: "/logo.svg",
      textLink: {
        text: "ignite.com",
        url: "https://ignite.com/",
      },
        {
          service: "twitter",
          url: "https://twitter.com/ignite_dev",
        },
        {
          service: "linkedin",
          url: "https://www.linkedin.com/company/ignt/",
        },
        {
          service: "discord",
          url: "https://discord.gg/ignt",
        },
        {
          service: "youtube",
          url: "https://www.youtube.com/ignitehq",
        },
      ],

      smallprint:
        "This website is maintained by Ignite. The contents and opinions of this website are those of Ignite.",
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
              url: "https://discord.gg/ignt",
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
