<template>
  <div class="table__panel panel">
    <div class="panel__top">
      <p>Transactions</p>
    </div>
    <div class="panel__table">
      <TableWrapper :tableHeads="['TxHash', 'Fee', 'Gas', 'Msgs']">

        <Accordion :id="'accordion-subTable-'+parentGroupId">
          <TableRowWrapper
            v-for="msg in messagesForTable"
            :key="msg.tableData.id"                          
          >
            <AccordionItem
              :itemData="msg.tableData"
              :groupId="parentGroupId"
            >            
              <TableRowCellsGroup      
                slot="trigger"       
                :tableCells="msg.txMsg"
              />                      
              <div slot="contents">
                <SideTabList
                  :list="getFmtMsgForInnerTable(msg.msgs)"
                />
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
import SideTabList from '@/components/table/SideTabList'

import Accordion from '@/components/accordion/Accordion'
import AccordionItem from '@/components/accordion/AccordionItem'

export default {
  components: {
    TableWrapper,
    TableRowWrapper,    
    TableRowCellsGroup,
    SideTabList,
    Accordion,
    AccordionItem
  },
  props: {
    parentGroupId: { type: String, require: true },
    // tableCells: { type: Array, require: true },
    rowItems: { type: Array, require: true }
  },
  data() {
    return {
      exampleData: [
        { id: 1, isActive: false },
        { id: 2, isActive: false }
      ]
    }
  },
  computed: {
    messagesForTable() {
      return this.rowItems.map(item => {
        const {
          fee,
          msg
        } = item

        return {
          txMsg: [
            'fakehashtestingssdf', // temp
            fee.amount[0].amount, // temp
            fee.gas, // temp
            msg.length
          ],
          msgs: msg.map(({
            type,
            value
          }) => ({
            type: this.getMsgType(type),
            amount: this.getAmount(value.amount),
            delegator: value.delegator_address,
            validator: value.validator_address,
            from: value.from_address,
            to: value.to_address
          })),
          tableData: {
            id: item.signatures[0].signature, // temp
            isActive: false
          },
        }
      })
    }
  },
  methods: {
    getAmount(amountObj) {
      return amountObj.amount 
        ? amountObj.amount+amountObj.denom
        : amountObj[0].amount+amountObj[0].denom
    },
    getFmtMsgForInnerTable(msgs) {
      return msgs.map(({
        type,
        amount,
        delegator,
        validator,
        from,
        to
      }) => {
        const fromType = type === 'MsgSend' ? {
          title: 'From', content: from
        } : {
          title: 'Delegator', content: delegator
        }
        const toType = type === 'MsgSend' ? {
          title: 'To', content: to
        } : {
          title: 'Validator', content: validator
        }

        return {
          title: type,
          subItems: [
            fromType,
            toType,
            { title: 'Amount', content: amount },
          ]
        }
      })
    },
    getMsgType(type) {
      return type.replace('cosmos-sdk/', '')
    }
  }
}
</script>

<style scoped>

.panel {
  background-color: var(--c-theme-secondary);
}

.panel.table__panel {
  padding-top: 1rem;
}

.panel .panel__top {
  padding-left: 1rem;
  padding-right: 1rem; 
}

.panel .panel__top p {
  margin: 0 0 1rem 0;
}

.panel >>> .table__row {
  padding-left: 1rem;
  padding-right: 1rem;
}
.panel >>> .table__row:last-child {
  padding-bottom: 0.5rem;
}
.panel >>> .table__row.-is-active {
  /* background-color: var(--); */
  background-color: var(--c-grey-3);
}
.panel >>> .table__row.-is-active:last-child {
  padding-bottom: 0;
}

.panel >>> .table__cells.-header {
  padding-left: 1rem;
  padding-right: 1rem;  
}
.panel >>> .table__cells {
  padding-top: 0.5rem;
  padding-bottom: 0.5rem;
  padding-left: 0;
  padding-right: 0;  
}

/* temporary table styling */
.panel__table >>> .table__cells .table__col:nth-child(1) {
  flex-grow: 1;
  width: auto;
}
.panel__table >>> .table__cells .table__col:nth-child(2) {
  width: 15%;
}
.panel__table >>> .table__cells .table__col:nth-child(3) {
  width: 15%;
}
.panel__table >>> .table__cells .table__col:nth-child(4) {
  flex-grow: 0;
  width: 15%;
}

</style>