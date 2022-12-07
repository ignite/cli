import React from "react";
import styled from "styled-components";

const GridItem = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  box-shadow: 1px 1px #efefef;
  font-weight: 500;
  padding: 1rem;
  text-align: center;
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  border-collapse: collapse;
  ${GridItem}:nth-child(3n) {
    box-shadow: 0 1px #efefef;
  }
  ${GridItem}:nth-last-child(1) {
    box-shadow: 1px 0px #efefef;
  }
  @media (max-width: 500px) {
    grid-template-columns: 1fr 1fr;
    ${GridItem}:nth-child(3n) {
      box-shadow: 1px 1px #efefef;
    }
    ${GridItem}:nth-child(2n) {
      box-shadow: 0 1px #efefef;
    }
    ${GridItem}:nth-last-child(1) {
      box-shadow: 0px 0px #efefef;
    }
    ${GridItem}:nth-last-child(2) {
      box-shadow: 1px 0px #efefef;
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
