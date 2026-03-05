// Copyright 2019 Huawei Technologies Co.,Ltd.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License.  You may obtain a copy of
// License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package obs

// Test fixtures for OBS SDK unit tests

const (
	// TestAK is the test access key ID
	TestAK = "test-access-key-id"
	// TestSK is the test secret access key
	TestSK = "test-secret-access-key"
	// TestEndpoint is the test endpoint URL
	TestEndpoint = "https://obs.cn-north-4.myhuaweicloud.com"
	// TestBucket is the test bucket name
	TestBucket = "test-bucket"
	// TestObjectKey is the test object key
	TestObjectKey = "test-object.txt"
	// TestSecurityToken is the test security token
	TestSecurityToken = "test-security-token"
	// TestRegion is the test region
	TestRegion = "cn-north-4"
)

// TestObjectContent is sample content for object tests
var TestObjectContent = []byte("This is test object content for unit tests.")

// TestObjectMetadata is sample metadata for object tests
var TestObjectMetadata = map[string]string{
	"Content-Type":  "text/plain",
	"Cache-Control":  "no-cache",
	"X-Test-Header": "test-value",
}

// TestBucketACL is sample ACL configuration for bucket tests
var TestBucketACL = &AccessControlPolicy{
	Owner: Owner{
		ID:          "test-owner-id",
		DisplayName: "test-owner",
	},
	Grants: []Grant{
		{
			Grantee: Grantee{
				Type: GranteeGroup,
				URI:  GroupAllUsers,
			},
			Permission: PermissionRead,
		},
	},
}

// TestObjectACL is sample ACL configuration for object tests
var TestObjectACL = &AccessControlPolicy{
	Owner: Owner{
		ID:          "test-object-owner-id",
		DisplayName: "test-object-owner",
	},
	Grants: []Grant{
		{
			Grantee: Grantee{
				Type: GranteeGroup,
				URI:  GroupAllUsers,
			},
			Permission: PermissionRead,
		},
		{
			Grantee: Grantee{
				Type: GranteeGroup,
				URI:  GroupAuthenticatedUsers,
			},
			Permission: PermissionFullControl,
		},
	},
}

// TestListBucketsXML is a sample ListBuckets response XML
const TestListBucketsXML = `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult>
	<Owner>
		<ID>test-owner-id</ID>
		<DisplayName>test-owner</DisplayName>
	</Owner>
	<Buckets>
		<Bucket>
			<Name>bucket1</Name>
			<CreationDate>2023-01-01T00:00:00Z</CreationDate>
			<Location>cn-north-4</Location>
		</Bucket>
		<Bucket>
			<Name>bucket2</Name>
			<CreationDate>2023-01-02T00:00:00Z</CreationDate>
			<Location>us-east-1</Location>
		</Bucket>
	</Buckets>
</ListAllMyBucketsResult>`

// TestListObjectsXML is a sample ListObjects response XML
const TestListObjectsXML = `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult>
	<Name>test-bucket</Name>
	<Prefix></Prefix>
	<Marker></Marker>
	<MaxKeys>1000</MaxKeys>
	<IsTruncated>false</IsTruncated>
	<Contents>
		<Key>object1.txt</Key>
		<LastModified>2023-01-01T00:00:00Z</LastModified>
		<ETag>"d41d8cd98f00b204e9800998ecf8427e"</ETag>
		<Size>1024</Size>
		<Owner>
			<ID>owner-id</ID>
			<DisplayName>owner</DisplayName>
		</Owner>
		<StorageClass>STANDARD</StorageClass>
	</Contents>
	<Contents>
		<Key>object2.jpg</Key>
		<LastModified>2023-01-02T00:00:00Z</LastModified>
		<ETag>"5d41402abc4b2a76b9719d911017c592"</ETag>
		<Size>2048</Size>
		<Owner>
			<ID>owner-id</ID>
			<DisplayName>owner</DisplayName>
		</Owner>
		<StorageClass>STANDARD</StorageClass>
	</Contents>
</ListBucketResult>`

// TestCreateBucketXML is a sample CreateBucket request XML
const TestCreateBucketXML = `<?xml version="1.0" encoding="UTF-8"?>
<CreateBucketConfiguration>
	<Location>cn-north-4</Location>
	<StorageClass>STANDARD</StorageClass>
</CreateBucketConfiguration>`

// TestErrorResponseXML is a sample error response XML
const TestErrorResponseXML = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
	<Code>NoSuchBucket</Code>
	<Message>The specified bucket does not exist</Message>
	<Resource>/test-bucket</Resource>
	<RequestId>test-request-id-123</RequestId>
	<HostId>test-host-id-456</HostId>
</Error>`

// TestAccessDeniedErrorXML is a sample AccessDenied error response XML
const TestAccessDeniedErrorXML = `<?xml version="1.0" encoding="UTF-8"?>
<Error>
	<Code>AccessDenied</Code>
	<Code>Access Denied</Code>
	<Message>Access Denied</Message>
	<Resource>/test-bucket/test-object.txt</Resource>
	<RequestId>test-request-id-789</RequestId>
	<HostId>test-host-id-012</HostId>
</Error>`

// TestListMultipartUploadsXML is a sample ListMultipartUploads response XML
const TestListMultipartUploadsXML = `<?xml version="1.0" encoding="UTF-8"?>
<ListMultipartUploadsResult>
	<Bucket>test-bucket</Bucket>
	<KeyMarker></KeyMarker>
	<UploadIdMarker></UploadIdMarker>
	<NextKeyMarker></NextKeyMarker>
	<NextUploadIdMarker></NextUploadIdMarker>
	<Delimiter></Delimiter>
	<Prefix></Prefix>
	<MaxUploads>1000</MaxUploads>
	<IsTruncated>false</IsTruncated>
	<Upload>
		<Key>test-large-file.zip</Key>
		<UploadId>upload-id-1234567890</UploadId>
		<Initiator>
			<ID>initiator-id</ID>
			<DisplayName>initiator</DisplayName>
		</Initiator>
		<Owner>
			<ID>owner-id</ID>
			<DisplayName>owner</DisplayName>
		</Owner>
		<StorageClass>STANDARD</StorageClass>
		<Initiated>2023-01-01T00:00:00Z</Initiated>
	</Upload>
</ListMultipartUploadsResult>`

// TestListPartsXML is a sample ListParts response XML
const TestListPartsXML = `<?xml version="1.0" encoding="UTF-8"?>
<ListPartsResult>
	<Bucket>test-bucket</Bucket>
	<Key>test-object.txt</Key>
	<UploadId>upload-id-1234567890</UploadId>
	<StorageClass>STANDARD</StorageClass>
	<PartNumberMarker>0</PartNumberMarker>
	<NextPartNumberMarker>2</NextPartNumberMarker>
	<MaxParts>1000</MaxParts>
	<IsTruncated>false</IsTruncated>
	<Part>
		<PartNumber>1</PartNumber>
		<LastModified>2023-01-01T00:00:00Z</LastModified>
		<ETag>"part1-etag"</ETag>
		<Size>5242880</Size>
	</Part>
	<Part>
		<PartNumber>2</PartNumber>
		<LastModified>2023-01-01T00:00:01Z</LastModified>
		<ETag>"part2-etag"</ETag>
		<Size>5242880</Size>
	</Part>
	<Initiator>
		<ID>initiator-id</ID>
		<DisplayName>initiator</DisplayName>
	</Initiator>
	<Owner>
		<ID>owner-id</ID>
		<DisplayName>owner</DisplayName>
	</Owner>
</ListPartsResult>`

// TestBucketPolicyXML is a sample bucket policy XML
const TestBucketPolicyXML = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Effect": "Allow",
			"Principal": {
				"AWS": ["*"]
			},
			"Action": ["s3:GetObject"],
			"Resource": ["arn:aws:s3:::test-bucket/*"]
		}
	]
}`

// TestBucketLifecycleXML is a sample bucket lifecycle configuration XML
const TestBucketLifecycleXML = `<?xml version="1.0" encoding="UTF-8"?>
<LifecycleConfiguration>
	<Rule>
		<ID>test-rule-id</ID>
		<Prefix>test-prefix/</Prefix>
		<Status>Enabled</Status>
		<Expiration>
			<Days>30</Days>
		</Expiration>
	</Rule>
</LifecycleConfiguration>`

// TestBucketCorsXML is a sample bucket CORS configuration XML
const TestBucketCorsXML = `<?xml version="1.0" encoding="UTF-8"?>
<CORSConfiguration>
	<CORSRule>
		<AllowedOrigin>*</AllowedOrigin>
		<AllowedMethod>GET</AllowedMethod>
		<AllowedMethod>PUT</AllowedMethod>
		<AllowedMethod>POST</AllowedMethod>
		<AllowedMethod>DELETE</AllowedMethod>
		<AllowedHeader>*</AllowedHeader>
		<MaxAgeSeconds>3000</MaxAgeSeconds>
		<ExposeHeader>ETag</ExposeHeader>
	</CORSRule>
</CORSConfiguration>`

// TestBucketWebsiteXML is a sample bucket website configuration XML
const TestBucketWebsiteXML = `<?xml version="1.0" encoding="UTF-8"?>
<WebsiteConfiguration>
	<IndexDocument>
		<Suffix>index.html</Suffix>
	</IndexDocument>
	<ErrorDocument>
		<Key>error.html</Key>
	</ErrorDocument>
</WebsiteConfiguration>`

// TestBucketVersioningXML is a sample bucket versioning configuration XML
const TestBucketVersioningXML = `<?xml version="1.0" encoding="UTF-8"?>
<VersioningConfiguration>
	<Status>Enabled</Status>
</VersioningConfiguration>`

// TestBucketACLXML is a sample bucket ACL response XML
const TestBucketACLXML = `<?xml version="1.0" encoding="UTF-8"?>
<AccessControlPolicy>
	<Owner>
		<ID>test-owner-id</ID>
		<DisplayName>test-owner</DisplayName>
	</Owner>
	<AccessControlList>
		<Grant>
			<Grantee>
				<ID>grantee-id</ID>
				<DisplayName>grantee</DisplayName>
			</Grantee>
			<Permission>READ</Permission>
		</Grant>
		<Grant>
			<Grantee>
				<Type>Group</Type>
				<URI>http://acs.amazonaws.com/groups/global/AllUsers</URI>
			</Grantee>
			<Permission>READ_ACP</Permission>
		</Grant>
	</AccessControlList>
</AccessControlPolicy>`

// TestObjectACLXML is a sample object ACL response XML
const TestObjectACLXML = `<?xml version="1.0" encoding="UTF-8"?>
<AccessControlPolicy>
	<Owner>
		<ID>test-object-owner-id</ID>
		<DisplayName>test-object-owner</DisplayName>
	</Owner>
	<AccessControlList>
		<Grant>
			<Grantee>
				<ID>grantee-id</ID>
				<DisplayName>grantee</DisplayName>
			</Grantee>
			<Permission>READ</Permission>
		</Grant>
	</AccessControlList>
</AccessControlPolicy>`
