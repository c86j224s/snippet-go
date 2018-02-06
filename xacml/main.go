package main

import (
	"encoding/xml"
	"fmt"
)

type AnyOf struct {
}

type AllOf struct {
}

type Target struct {
	AnyOf []AnyOf `xml:"AnyOf"`
	AllOf []AllOf `xml:"AllOf"`
}

type Rule struct {
	RuleID      string   `xml:"RuleId,attr"`
	Effect      string   `xml:"Effect,attr"`
	Description string   `xml:"Description"`
	Target      []Target `xml:"Target"`
}

type Policy struct {
	Xmlns              string   `xml:"xmlns,attr"`
	XmlnsXsi           string   `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation  string   `xml:"xsi:xchemaLocation,attr"`
	PolicyID           string   `xml:"PolicyId,attr"`
	Version            string   `xml:"Version,attr"`
	RuleCombiningAlgID string   `xml:"RuleCombiningAlgId,attr"`
	Description        string   `xml:"Description"`
	Target             []Target `xml:"Target"`
	Rule               []Rule   `xml:"Rule"`
}

/*
type PolicySet struct {
	Name xml.Name
	Policy []Policy
	Rule   []Rule
}
*/

func main() {
	data := `
<?xml version="1.0" encoding="UTF-8"?>
<Policy
	xmlns="urn:oasis:names:tc:xacml:3.0:core:schema:wd-17"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xsi:schemaLocation="urn:oasis:names:tc:xacml:3.0:core:schema:wd-17
		http://docs.oasis-open.org/xacml/3.0/xacml-core-v3-schema-wd-17.xsd"
	PolicyId="urn:oasis:names:tc:xacml:3.0:example:SimplePolicy1"
	Version="1.0"
	RuleCombiningAlgId="identifier:rule-combining-algorithm:deny-overrides">
	
	<Description>
		Medi Corp access control policy
	</Description>
	<Target/>
	<Rule
		RuleId= "urn:oasis:names:tc:xacml:3.0:example:SimpleRule1"
		Effect="Permit">

		<Description>
			Any subject with an e-mail name in the med.example.com domain
			can perform any action on any resource.
		</Description>
		<Target>
			<AnyOf>
				<AllOf>
					<Match
						MatchId="urn:oasis:names:tc:xacml:1.0:function:rfc822Name-match">
					
						<AttributeValue
							DataType="http://www.w3.org/2001/XMLSchema#string"
						>med.example.com</AttributeValue>
						<AttributeDesignator
							MustBePresent="false"
							Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
							AttributeId="urn:oasis:names:tc:xacml:1.0:subject:subject-id"
							DataType="urn:oasis:names:tc:xacml:1.0:data-type:rfc822Name"/>
					</Match>
				</AllOf>
			</AnyOf>
		</Target>
	</Rule>
</Policy>`

	v := Policy{}

	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		fmt.Printf("error :%v", err)
		return
	}
}

/*
	xml handling in go lang
	- https://golang.org/pkg/encoding/xml/#Unmarshal
	- http://golang.site/go/article/105-XML-%EC%82%AC%EC%9A%A9
	- https://www.joinc.co.kr/w/man/12/golang/networkProgramming/xml
	- https://golang.org/src/encoding/xml/example_test.go

	xacml
	- http://docs.oasis-open.org/xacml/3.0/xacml-3.0-core-spec-os-en.html
*/
