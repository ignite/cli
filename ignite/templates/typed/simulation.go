package typed

import (
	"fmt"
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/multiformatname"
	"github.com/ignite/cli/v29/ignite/pkg/placeholder"
)

func ModuleSimulationMsgModify(
	replacer placeholder.Replacer,
	content string,
	typeName, msgSigner multiformatname.Name,
	msgs ...string,
) string {
	if len(msgs) == 0 {
		msgs = append(msgs, "")
	}

	for _, msg := range msgs {
		// simulation operations
		templateOp := `reg.Add(weights.Get("msg_%[4]v", 100 /* determine the simulation weight value */), simulation.Msg%[2]v%[3]vFactory(am.keeper))
	%[1]v`
		replacementOp := fmt.Sprintf(templateOp, PlaceholderSimappOperation, msg, typeName.UpperCamel, fmt.Sprintf("%s_%s", strings.ToLower(msg), typeName.Snake))
		content = replacer.Replace(content, PlaceholderSimappOperation, replacementOp)

		// add proposal simulation operations for msgs having an authority as signer.
		if strings.Contains(content, PlaceholderSimappOperationProposalSimulation) && strings.EqualFold(msgSigner.Original, "authority") {
			templateOpMsg := `reg.Add(weights.Get("msg_%[4]v", 100), simulation.Msg%[2]v%[3]vFactory(am.keeper))
	%[1]v`
			replacementOpMsg := fmt.Sprintf(templateOpMsg, PlaceholderSimappOperationProposalSimulation, msg, typeName.UpperCamel, fmt.Sprintf("%s_%s", strings.ToLower(msg), typeName.Snake))
			content = replacer.Replace(content, PlaceholderSimappOperationProposalSimulation, replacementOpMsg)
		}
	}

	return content
}
