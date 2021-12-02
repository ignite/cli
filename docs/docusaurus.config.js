// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require("prism-react-renderer/themes/github");
const darkCodeTheme = require("prism-react-renderer/themes/dracula");

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: "Starport",
  tagline: "Starport",
  url: "https://docs.starport.com",
  baseUrl: "/",
  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/favicon.ico",
  organizationName: "tendermint", // Usually your GitHub org/user name.
  projectName: "starport", // Usually your repo name.
  presets: [
    [
      "@docusaurus/preset-classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          // Please change this to your repo.
          editUrl: "https://github.com/tendermint/starport/edit/develop/",
          routeBasePath: "/",
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        logo: {
          alt: "My Site Logo",
          src: "img/logo.svg",
          srcDark: "img/logo_dark.svg",
          width: 220,
        },
        items: [
          {
            href: "https://starport.com",
            label: "starport.com",
            position: "right",
          },
          {
            href: "https://github.com/tendermint/starport",
            label: "GitHub",
            position: "right",
          },
        ],
      },
      footer: {
        style: "dark",
        logo: {
          src: "img/logo_dark.svg",
        },
        links: [
          {
            items: [
              {
                label: "Starport website",
                href: "https://starport.com",
              },
              {
                label: "Try Starport online",
                href: "https://gitpod.io/#https://github.com/tendermint/starport/tree/master",
              },
              {
                label: "Blog",
                href: "https://starport.com/blog",
              },
            ],
          },
          {
            items: [
              {
                label: "Twitter",
                href: "https://twitter.com/StarportHQ",
              },
              {
                label: "Developer chat",
                href: "https://discord.gg/starport",
              },
            ],
          },
          {
            items: [
              {
                label: "GitHub",
                href: "https://github.com/tendermint/starport",
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Tendermint Inc.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
    }),
};

module.exports = config;
