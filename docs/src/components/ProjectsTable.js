import React from "react";
import styled from "styled-components";

const GridItem = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  box-shadow: 1px 1px var(--ifm-color-emphasis-200);
  font-weight: 500;
  padding: 1rem;
  text-align: center;
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  border-collapse: collapse;
  position: relative;
  &:after {
    width: 100%;
    height: 1px;
    content: "";
    bottom: -1px;
    background: var(--ifm-background-color);
    position: absolute;
  }
  ${GridItem}:nth-child(3n) {
    box-shadow: 0 1px var(--ifm-color-emphasis-200);
  }
  html[data-theme="dark"] & img {
    filter: invert(1);
  }
  @media (max-width: 500px) {
    grid-template-columns: 1fr 1fr;
    ${GridItem}:nth-child(3n) {
      box-shadow: 1px 1px var(--ifm-color-emphasis-200);
    }
    ${GridItem}:nth-child(2n) {
      box-shadow: 0 1px var(--ifm-color-emphasis-200);
    }
  }
`;

export default function ProjectsTable({ data }) {
  return (
    <Grid>
      {data.map((item) => (
        <GridItem key={item.logo}>
          <img src={item.logo} />
          <div>{item.name}</div>
        </GridItem>
      ))}
    </Grid>
  );
}
