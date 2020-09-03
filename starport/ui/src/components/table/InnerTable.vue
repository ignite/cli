<template>
  <div class="table__panel panel">
    <div class="panel__top">
      <p>Transactions</p>
    </div>
    <div class="panel__table">
      <TableWrapper :tableHeads="['TxHash', 'Result', 'Fee', 'Gas', 'Msgs']">

        <Accordion :id="'accordion-subTable-'+parentGroupId">
          <TableRowWrapper
            v-for="item in exampleData"
            :key="item.id"                          
          >
            <AccordionItem
              :itemData="item"
              :groupId="parentGroupId"
            >            
              <TableRowCellsGroup      
                slot="trigger"       
                :tableCells="['16s ago', '3111975', 'stake.fish', 'BBE79CF80CD...70A9F9F436A', '3']"
              />                      
              <div slot="contents">
                <p>testing 123</p>
              </div>                    
            </AccordionItem>
          </TableRowWrapper>
        </Accordion>

      </TableWrapper>
    </div>
  </div>      
</template>

<script>
import TableWrapper from '@/components/table/TableWrapper'
import TableRowWrapper from '@/components/table/RowWrapper'
import TableRowCellsGroup from '@/components/table/RowCellsGroup'

import Accordion from '@/components/accordion/Accordion'
import AccordionItem from '@/components/accordion/AccordionItem'

export default {
  components: {
    TableWrapper,
    TableRowWrapper,    
    TableRowCellsGroup,
    Accordion,
    AccordionItem
  },
  props: [
    'parentGroupId'
  ],
  data() {
    return {
      exampleData: [
        { id: 1, isActive: false },
        { id: 2, isActive: false }
      ]
    }
  }  
}
</script>

<style scoped>

</style>