package protoanalysis

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/emicklei/proto"
)

type builder struct {
	p protoPackage
}

// build turns a low level proto pkg into a high level Package.
func build(p protoPackage) Package {
	br := builder{p}

	pk := Package{
		Name:     p.name,
		Path:     p.dir,
		Files:    br.buildFiles(),
		Messages: br.buildMessages(),
		Services: br.toServices(p.services()),
	}

	for _, option := range p.options() {
		if option.Name == optionGoPkg {
			pk.GoImportName = option.Constant.Source
			break
		}
	}

	return pk
}

func (b builder) buildFiles() (files []File) {
	for _, f := range b.p.files {
		files = append(files, File{f.path, f.imports})
	}

	return
}

func (b builder) buildMessages() (messages []Message) {
	for _, f := range b.p.files {
		for _, message := range f.messages {
			// Keep track of the message fields and types
			fields := make(map[string]string)

			// Find the highest field number
			var highestFieldNumber int
			for _, elem := range message.Elements {
				field, ok := elem.(*proto.NormalField)
				if !ok {
					continue
				}

				if field.Sequence > highestFieldNumber {
					highestFieldNumber = field.Sequence
				}

				fields[field.Name] = field.Type
			}

			// some proto messages might be defined inside another proto messages.
			// to represents these types, an underscore is used.
			// e.g. if C message inside B, and B inside A: A_B_C.
			var (
				name   = message.Name
				parent = message.Parent
			)
			for {
				if parent == nil {
					break
				}

				parentMessage, ok := parent.(*proto.Message)
				if !ok {
					break
				}

				name = fmt.Sprintf("%s_%s", parentMessage.Name, name)
				parent = parentMessage.Parent
			}

			messages = append(messages, Message{
				Name:               name,
				Path:               f.path,
				HighestFieldNumber: highestFieldNumber,
				Fields:             fields,
			})
		}
	}

	return messages
}

func (b builder) toServices(ps []*proto.Service) (services []Service) {
	for _, service := range ps {
		s := Service{
			Name:     service.Name,
			RPCFuncs: b.elementsToRPCFunc(service.Elements),
		}

		services = append(services, s)
	}

	return
}

func (b builder) elementsToRPCFunc(elems []proto.Visitee) (rpcFuncs []RPCFunc) {
	for _, el := range elems {
		rpc, ok := el.(*proto.RPC)
		if !ok {
			continue
		}

		var requestMessage *proto.Message

		for _, message := range b.p.messages() {
			if message.Name != rpc.RequestType {
				continue
			}
			requestMessage = message
		}

		if requestMessage == nil {
			continue
		}

		rf := RPCFunc{
			Name:        rpc.Name,
			RequestType: rpc.RequestType,
			ReturnsType: rpc.ReturnsType,
			HTTPRules:   b.elementsToHTTPRules(requestMessage, rpc.Elements),
		}

		rpcFuncs = append(rpcFuncs, rf)
	}

	return rpcFuncs
}

func (b builder) elementsToHTTPRules(requestMessage *proto.Message, elems []proto.Visitee) (httpRules []HTTPRule) {
	for _, el := range elems {
		option, ok := el.(*proto.Option)
		if !ok {
			continue
		}
		if !strings.Contains(option.Name, "google.api.http") {
			continue
		}

		httpRules = append(httpRules, b.constantToHTTPRules(requestMessage, option.Constant)...)
	}

	return
}

// Regexp to extract HTTP rule URL parameter names.
// The expression extracts parameter names defined within "{}".
// Extra parameter arguments are ignored. These arguments are normally
// defined after an "=", for example as "{param=**}".
var urlParamRe = regexp.MustCompile(`(?m){([^=]+?)(?:=.+?)?}`)

func (b builder) constantToHTTPRules(requestMessage *proto.Message, constant proto.Literal) (httpRules []HTTPRule) {
	// find out the endpoint template.
	endpoint := constant.Source

	if endpoint == "" {
		for _, each := range constant.OrderedMap {
			switch each.Name {
			case
				"get",
				"post",
				"put",
				"patch",
				"delete":
				endpoint = each.Source
			}
			if endpoint != "" {
				break
			}
		}
	}

	// find out url params.
	var params []string

	match := urlParamRe.FindAllStringSubmatch(endpoint, -1)
	for _, item := range match {
		params = append(params, item[1])
	}

	// calculate url params, query params and body fields counts.
	var (
		messageFields, messageFieldsCount = b.messageFieldsCount(requestMessage)
		paramsCount                       = len(params)
		bodyFieldsCount                   int
	)

	if body, ok := constant.OrderedMap.Get("body"); ok { // check if body is specified.
		if body.Source == "*" { // means there should be no query params per the spec.
			bodyFieldsCount = messageFieldsCount - paramsCount
		} else if body.Source != "" {
			bodyFieldsCount = 1 // means body fields are grouped under a single top-level field.
		}
	}

	queryParamsCount := messageFieldsCount - paramsCount - bodyFieldsCount

	var (
		queryFields map[string]string
		bodyFields  map[string]string
	)
	for name, t := range messageFields {
		if slices.Contains(params, name) {
			// this is a URL parameter, skip it
			continue
		}

		// If there are body fields, we need to add them to the bodyFields map.
		// There are no known post requests that contain body fields and query params
		if bodyFieldsCount > 0 {
			if len(bodyFields) == 0 {
				bodyFields = make(map[string]string)
			}
			bodyFields[name] = t
		} else {
			if len(queryFields) == 0 {
				queryFields = make(map[string]string)
			}

			queryFields[name] = t
		}
	}

	// create and add the HTTP rule to the list.
	httpRule := HTTPRule{
		Endpoint:    endpoint,
		Params:      params,
		HasQuery:    queryParamsCount > 0,
		QueryFields: queryFields,
		HasBody:     bodyFieldsCount > 0,
		BodyFields:  bodyFields,
	}

	httpRules = append(httpRules, httpRule)

	// search for nested HTTP rules.
	if constant, ok := constant.OrderedMap.Get("additional_bindings"); ok {
		httpRules = append(httpRules, b.constantToHTTPRules(requestMessage, *constant)...)
	}

	return httpRules
}

func (b builder) messageFieldsCount(message *proto.Message) (messageFields map[string]string, count int) {
	messageFields = make(map[string]string)

	for _, el := range message.Elements {
		switch el := el.(type) {
		case *proto.NormalField:
			count++
			if el.Repeated {
				messageFields[el.Name] = fmt.Sprintf("repeated %s", el.Type)
			} else {
				messageFields[el.Name] = el.Type
			}
		case *proto.MapField:
			count++
			messageFields[el.Name] = fmt.Sprintf("map<%s, %s>", el.KeyType, el.Type)
		case *proto.OneOfField:
			count++
			messageFields[el.Name] = el.Type
		}
	}

	return
}
