package xmlutil_test

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jviney/aws-sdk-go-v2/aws"
	"github.com/jviney/aws-sdk-go-v2/private/protocol/xml/xmlutil"
	"github.com/jviney/aws-sdk-go-v2/service/ec2"
)

type DataOutput struct {
	_ struct{} `type:"structure"`

	FooEnum string `type:"string" enum:"true"`

	ListEnums []string `type:"list"`
}

func BenchmarkXMLUnmarshal_Simple(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := getXMLResponseSimple()
		xmlutil.UnmarshalXML(req.Data, xml.NewDecoder(req.HTTPResponse.Body), "")
	}
}

func BenchmarkXMLUnmarshal_Complex(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := getXMLResponseComplex()
		xmlutil.UnmarshalXML(req.Data, xml.NewDecoder(req.HTTPResponse.Body), "")
	}
}

func getXMLResponseSimple() *aws.Request {
	buf := bytes.NewReader([]byte("<OperationNameResponse><FooEnum>foo</FooEnum><ListEnums><member>0</member><member>1</member></ListEnums></OperationNameResponse>"))
	req := aws.Request{Data: &DataOutput{}, HTTPResponse: &http.Response{Body: ioutil.NopCloser(buf)}}
	return &req
}

func getXMLResponseComplex() *aws.Request {
	buf := bytes.NewReader([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<DescribeInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2016-11-15/">
    <requestId>7e2ca54c-e2af-4567-bb41-21632d2b839e</requestId>
    <reservationSet>
        <item>
            <reservationId>r-05e953f164a34c484</reservationId>
            <ownerId>183557167593</ownerId>
            <groupSet/>
            <instancesSet>
                <item>
                    <instanceId>i-05805668ced0206f0</instanceId>
                    <imageId>ami-082b5a644766e0e6f</imageId>
                    <instanceState>
                        <code>16</code>
                        <name>running</name>
                    </instanceState>
                    <privateDnsName>ip-172-31-30-42.us-west-2.compute.internal</privateDnsName>
                    <dnsName>ec2-34-219-17-124.us-west-2.compute.amazonaws.com</dnsName>
                    <reason/>
                    <keyName>ec2</keyName>
                    <amiLaunchIndex>0</amiLaunchIndex>
                    <productCodes/>
                    <instanceType>t2.micro</instanceType>
                    <launchTime>2019-07-01T21:15:47.000Z</launchTime>
                    <placement>
                        <availabilityZone>us-west-2a</availabilityZone>
                        <groupName/>
                        <tenancy>default</tenancy>
                    </placement>
                    <monitoring>
                        <state>disabled</state>
                    </monitoring>
                    <subnetId>subnet-21959558</subnetId>
                    <vpcId>vpc-1de55365</vpcId>
                    <privateIpAddress>172.31.30.42</privateIpAddress>
                    <ipAddress>34.219.17.124</ipAddress>
                    <sourceDestCheck>true</sourceDestCheck>
                    <groupSet>
                        <item>
                            <groupId>sg-02d1f51eb2fa52795</groupId>
                            <groupName>launch-wizard-2</groupName>
                        </item>
                    </groupSet>
                    <architecture>x86_64</architecture>
                    <rootDeviceType>ebs</rootDeviceType>
                    <rootDeviceName>/dev/xvda</rootDeviceName>
                    <blockDeviceMapping>
                        <item>
                            <deviceName>/dev/xvda</deviceName>
                            <ebs>
                                <volumeId>vol-08225e4fc2fde8e73</volumeId>
                                <status>attached</status>
                                <attachTime>2019-07-01T21:15:48.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </ebs>
                        </item>
                    </blockDeviceMapping>
                    <virtualizationType>hvm</virtualizationType>
                    <clientToken/>
                    <hypervisor>xen</hypervisor>
                    <networkInterfaceSet>
                        <item>
                            <networkInterfaceId>eni-0ba368b59d3f5230e</networkInterfaceId>
                            <subnetId>subnet-21959558</subnetId>
                            <vpcId>vpc-1de55365</vpcId>
                            <description/>
                            <ownerId>183557167593</ownerId>
                            <status>in-use</status>
                            <macAddress>02:36:86:6e:84:7c</macAddress>
                            <privateIpAddress>172.31.30.42</privateIpAddress>
                            <privateDnsName>ip-172-31-30-42.us-west-2.compute.internal</privateDnsName>
                            <sourceDestCheck>true</sourceDestCheck>
                            <groupSet>
                                <item>
                                    <groupId>sg-02d1f51eb2fa52795</groupId>
                                    <groupName>launch-wizard-2</groupName>
                                </item>
                            </groupSet>
                            <attachment>
                                <attachmentId>eni-attach-0d52b5e24dcb77ede</attachmentId>
                                <deviceIndex>0</deviceIndex>
                                <status>attached</status>
                                <attachTime>2019-07-01T21:15:47.000Z</attachTime>
                                <deleteOnTermination>true</deleteOnTermination>
                            </attachment>
                            <association>
                                <publicIp>34.219.17.124</publicIp>
                                <publicDnsName>ec2-34-219-17-124.us-west-2.compute.amazonaws.com</publicDnsName>
                                <ipOwnerId>amazon</ipOwnerId>
                            </association>
                            <privateIpAddressesSet>
                                <item>
                                    <privateIpAddress>172.31.30.42</privateIpAddress>
                                    <privateDnsName>ip-172-31-30-42.us-west-2.compute.internal</privateDnsName>
                                    <primary>true</primary>
                                    <association>
                                    <publicIp>34.219.17.124</publicIp>
                                    <publicDnsName>ec2-34-219-17-124.us-west-2.compute.amazonaws.com</publicDnsName>
                                    <ipOwnerId>amazon</ipOwnerId>
                                    </association>
                                </item>
                            </privateIpAddressesSet>
                            <ipv6AddressesSet/>
                            <interfaceType>interface</interfaceType>
                        </item>
                    </networkInterfaceSet>
                    <ebsOptimized>false</ebsOptimized>
                    <enaSupport>true</enaSupport>
                    <cpuOptions>
                        <coreCount>1</coreCount>
                        <threadsPerCore>1</threadsPerCore>
                    </cpuOptions>
                    <capacityReservationSpecification>
                        <capacityReservationPreference>open</capacityReservationPreference>
                    </capacityReservationSpecification>
                    <hibernationOptions>
                        <configured>false</configured>
                    </hibernationOptions>
                </item>
            </instancesSet>
        </item>
    </reservationSet>
</DescribeInstancesResponse>`))
	req := aws.Request{Data: &ec2.DescribeInstancesOutput{}, HTTPResponse: &http.Response{Body: ioutil.NopCloser(buf)}}
	return &req
}
