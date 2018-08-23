package BLC

import (
	"fmt"
	"os"
)

func (cli *CLIYS) sendYS(from []string, to []string, amount []string) {
	/*
	a:  12Gzy87R7dS19bf6VSfqG1qMJoWSUWNRDoThCtgEF5WdLMxJrWp
	b:  12qg59rrjzA1DFGUTXB5Cbbkh244RtJydwQDhGcvrJpP5VcfwTP
	c:  123SEF6i4vxhMcYyrQ1fckTsqq3oMYrqyZZ1eX4CPzDy57eYjuQ

	go run main.go send -from '["12Gzy87R7dS19bf6VSfqG1qMJoWSUWNRDoThCtgEF5WdLMxJrWp"]' -to '["12qg59rrjzA1DFGUTXB5Cbbkh244RtJydwQDhGcvrJpP5VcfwTP"]' -amount '["4"]'
	go run main.go send -from '["12Gzy87R7dS19bf6VSfqG1qMJoWSUWNRDoThCtgEF5WdLMxJrWp","12Gzy87R7dS19bf6VSfqG1qMJoWSUWNRDoThCtgEF5WdLMxJrWp"]' -to '["12qg59rrjzA1DFGUTXB5Cbbkh244RtJydwQDhGcvrJpP5VcfwTP","123SEF6i4vxhMcYyrQ1fckTsqq3oMYrqyZZ1eX4CPzDy57eYjuQ"]' -amount '["2","1"]'
	go run main.go send -from '["12qg59rrjzA1DFGUTXB5Cbbkh244RtJydwQDhGcvrJpP5VcfwTP","123SEF6i4vxhMcYyrQ1fckTsqq3oMYrqyZZ1eX4CPzDy57eYjuQ"]' -to '["123SEF6i4vxhMcYyrQ1fckTsqq3oMYrqyZZ1eX4CPzDy57eYjuQ","12Gzy87R7dS19bf6VSfqG1qMJoWSUWNRDoThCtgEF5WdLMxJrWp"]' -amount '["3","1"]'
	go run main.go send -from '["12Gzy87R7dS19bf6VSfqG1qMJoWSUWNRDoThCtgEF5WdLMxJrWp"]' -to '["123SEF6i4vxhMcYyrQ1fckTsqq3oMYrqyZZ1eX4CPzDy57eYjuQ"]' -amount '["8"]'


	1/	a->b 4					a: 7 / b: 4 / c: 0
	2/	a->b 2  a->c 1			a: 5 / b: 6 / c: 1
	3/	b->c 3  c->a 1			a: 6 / b: 4 / c: 3
	4/  a->c 8					a: 6 / b: 4 / c: 3
	 */

	if DBExistsYS() == false {
		fmt.Println("数据不存在.......")
		os.Exit(1)
	}

	blockchain := BlockchainObjectYS()
	defer blockchain.DBYS.Close()

	blockchain.MineNewBlockYS(from, to, amount)

	utxoSet := &UTXOSetYS{blockchain}
	utxoSet.UpdateYS()
}
