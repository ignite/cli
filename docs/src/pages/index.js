import React from "react";
import Layout from "@theme/Layout";
import Link from "@docusaurus/Link";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";
import styles from "./index.module.css";

const HomeCard = ({ title, description, to }) => (
  <div className={styles.card}>
    <h2>{title}</h2>
    <p>{description}</p>
    <Link className={styles.cardLink} to={to}>
      Explore →
    </Link>
  </div>
);

export default function Home() {
  const { siteConfig } = useDocusaurusContext();

  return (
    <Layout
      title={`${siteConfig.title} - Welcome`}
      description="Welcome to IGNITE® Documentation Portal"
    >
      <div className={styles.hero}>
        <h1 className={styles.heroTitle}>Welcome to IGNITE® Knowledge Hub</h1>
        <p className={styles.heroSubtitle}>
          Your one-stop portal for IGNITE® documentation, tutorials, and
          resources
        </p>
      </div>

      <main className={styles.main}>
        <div className={styles.cardsContainer}>
          <HomeCard
            title="IGNITE® Docs"
            description="Comprehensive documentation for IGNITE® CLI"
            to="/welcome"
          />
          <HomeCard
            title="IGNITE® Tutorials"
            description="Step-by-step guides and learning resources"
            to="https://tutorials.ignite.com"
          />
          <HomeCard
            title="IGNITE® Apps"
            description="Explore recommended IGNITE® Apps"
            to="https://ignite.com/marketplace"
          />
          <HomeCard
            title="Community"
            description="Join the IGNITE® community and connect with others"
            to="https://discord.com/invite/ignitecli"
          />
        </div>
      </main>
    </Layout>
  );
}
