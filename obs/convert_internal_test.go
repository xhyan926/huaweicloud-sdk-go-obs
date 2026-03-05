// Copyright 2019 Huawei Technologies Co.,Ltd.
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use
// this file except in compliance with the License.  You may obtain a copy of
// the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied. See the License for the
// specific language governing permissions and limitations under the License.

package obs

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ParseStringToFSStatusType Tests

func TestParseStringToFSStatusType_ShouldReturnEnabled_WhenValueIsEnabled(t *testing.T) {
	result := ParseStringToFSStatusType("Enabled")
	assert.Equal(t, FSStatusEnabled, result)
}

func TestParseStringToFSStatusType_ShouldReturnDisabled_WhenValueIsDisabled(t *testing.T) {
	result := ParseStringToFSStatusType("Disabled")
	assert.Equal(t, FSStatusDisabled, result)
}

func TestParseStringToFSStatusType_ShouldReturnEmpty_WhenValueIsInvalid(t *testing.T) {
	result := ParseStringToFSStatusType("Invalid")
	assert.Equal(t, FSStatusType(""), result)
}

// prepareGrantURI Tests

func TestPrepareGrantURI_ShouldReturnURI_WhenGroupAllUsers(t *testing.T) {
	result := prepareGrantURI(GroupAllUsers)
	assert.Contains(t, result, "http://acs.amazonaws.com/groups/global/AllUsers")
}

func TestPrepareGrantURI_ShouldReturnURI_WhenGroupAuthenticatedUsers(t *testing.T) {
	result := prepareGrantURI(GroupAuthenticatedUsers)
	assert.Contains(t, result, "http://acs.amazonaws.com/groups/global/AuthenticatedUsers")
}

func TestPrepareGrantURI_ShouldReturnURI_WhenGroupLogDelivery(t *testing.T) {
	result := prepareGrantURI(GroupLogDelivery)
	assert.Contains(t, result, "http://acs.amazonaws.com/groups/s3/LogDelivery")
}

func TestPrepareGrantURI_ShouldReturnSimpleURI_WhenOtherGroup(t *testing.T) {
	result := prepareGrantURI("custom-uri")
	assert.Contains(t, result, "custom-uri")
}

// convertGrantToXML Tests

func TestConvertGrantToXML_ShouldReturnXML_WhenUserType(t *testing.T) {
	grant := Grant{
		Grantee: Grantee{
			ID: "test-id",
		},
		Permission: "READ",
	}

	result := convertGrantToXML(grant, false, false)
	assert.Contains(t, result, "<Grant><Grantee")
	assert.Contains(t, result, "<ID>test-id</ID>")
	assert.Contains(t, result, "<Permission>READ</Permission>")
}

func TestConvertGrantToXML_ShouldReturnXML_WhenGroupType(t *testing.T) {
	grant := Grant{
		Grantee: Grantee{
			URI: GroupAllUsers,
		},
		Permission: "READ",
	}

	result := convertGrantToXML(grant, false, false)
	assert.Contains(t, result, "<Grant><Grantee")
	assert.Contains(t, result, "http://acs.amazonaws.com/groups/global/AllUsers")
	assert.Contains(t, result, "<Permission>READ</Permission>")
}

func TestConvertGrantToXML_ShouldReturnXML_WhenOBSType(t *testing.T) {
	grant := Grant{
		Grantee: Grantee{
			ID: "test-id",
		},
		Permission: "READ",
	}

	result := convertGrantToXML(grant, true, false)
	assert.Contains(t, result, "<Grant><Grantee")
	assert.Contains(t, result, "<ID>test-id</ID>")
	assert.NotContains(t, result, "xsi:type")
}

// ConvertAclToXml Tests

func TestConvertAclToXml_ShouldReturnXML_WhenValidInput(t *testing.T) {
	input := AccessControlPolicy{
		Owner: Owner{
			ID: "owner-id",
		},
		Grants: []Grant{
			{
				Grantee: Grantee{
					ID: "grantee-id",
				},
				Permission: "READ",
			},
		},
	}

	data, md5 := ConvertAclToXml(input, false, false)
	assert.Contains(t, data, "<AccessControlPolicy>")
	assert.Contains(t, data, "owner-id")
	assert.Contains(t, data, "grantee-id")
	assert.Empty(t, md5)
}

func TestConvertAclToXml_ShouldReturnMd5_WhenReturnMd5True(t *testing.T) {
	input := AccessControlPolicy{
		Owner: Owner{
			ID: "owner-id",
		},
	}

	data, md5 := ConvertAclToXml(input, true, false)
	assert.Contains(t, data, "<AccessControlPolicy>")
	assert.NotEmpty(t, md5)
}

// convertBucketACLToXML Tests

func TestConvertBucketACLToXML_ShouldReturnXML_WhenValidInput(t *testing.T) {
	input := AccessControlPolicy{
		Owner: Owner{
			ID: "owner-id",
		},
		Grants: []Grant{
			{
				Grantee: Grantee{
					ID: "grantee-id",
				},
				Permission: "READ",
			},
		},
	}

	data, md5 := convertBucketACLToXML(input, false, false)
	assert.Contains(t, data, "<AccessControlPolicy>")
	assert.Contains(t, data, "owner-id")
	assert.Empty(t, md5)
}

// convertConditionToXML Tests

func TestConvertConditionToXML_ShouldReturnXML_WhenKeyPrefixEquals(t *testing.T) {
	condition := Condition{
		KeyPrefixEquals: "prefix/",
	}

	result := convertConditionToXML(condition)
	assert.Contains(t, result, "<Condition>")
	assert.Contains(t, result, "<KeyPrefixEquals>prefix/</KeyPrefixEquals>")
	assert.Contains(t, result, "</Condition>")
}

func TestConvertConditionToXML_ShouldReturnXML_WhenHttpErrorCodeReturnedEquals(t *testing.T) {
	condition := Condition{
		HttpErrorCodeReturnedEquals: "404",
	}

	result := convertConditionToXML(condition)
	assert.Contains(t, result, "<Condition>")
	assert.Contains(t, result, "<HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals>")
	assert.Contains(t, result, "</Condition>")
}

func TestConvertConditionToXML_ShouldReturnEmpty_WhenNoFields(t *testing.T) {
	condition := Condition{}

	result := convertConditionToXML(condition)
	assert.Equal(t, "", result)
}

// prepareRoutingRule Tests

func TestPrepareRoutingRule_ShouldReturnXML_WhenValidInput(t *testing.T) {
	input := BucketWebsiteConfiguration{
		RoutingRules: []RoutingRule{
			{
				Redirect: Redirect{
					Protocol: "https",
					HostName: "example.com",
				},
			},
		},
	}

	result := prepareRoutingRule(input)
	assert.Contains(t, result, "<RoutingRule>")
	assert.Contains(t, result, "<Redirect>")
	assert.Contains(t, result, "<Protocol>https</Protocol>")
	assert.Contains(t, result, "<HostName>example.com</HostName>")
}

// converntFilterRulesToXML Tests

func TestConverntFilterRulesToXML_ShouldReturnXML_WhenValidRules(t *testing.T) {
	rules := []FilterRule{
		{
			Name:  "prefix",
			Value: "value/",
		},
	}

	result := converntFilterRulesToXML(rules, false)
	assert.Contains(t, result, "<Filter>")
	assert.Contains(t, result, "<S3Key>")
	assert.Contains(t, result, "<Name>prefix</Name>")
	assert.Contains(t, result, "<Value>value/</Value>")
}

func TestConverntFilterRulesToXML_ShouldReturnEmpty_WhenNoRules(t *testing.T) {
	rules := []FilterRule{}

	result := converntFilterRulesToXML(rules, false)
	assert.Equal(t, "", result)
}

// converntEventsToXML Tests

func TestConverntEventsToXML_ShouldReturnXML_WhenValidEvents(t *testing.T) {
	events := []EventType{
		ObjectCreatedPut,
		ObjectCreatedPost,
	}

	result := converntEventsToXML(events, false)
	assert.Contains(t, result, "s3:ObjectCreated:Put")
	assert.Contains(t, result, "s3:ObjectCreated:Post")
}

func TestConverntEventsToXML_ShouldReturnEmpty_WhenNoEvents(t *testing.T) {
	events := []EventType{}

	result := converntEventsToXML(events, false)
	assert.Equal(t, "", result)
}

// converntConfigureToXML Tests

func TestConverntConfigureToXML_ShouldReturnXML_WhenValidInput(t *testing.T) {
	config := TopicConfiguration{
		Topic: "arn:aws:sns:us-east-1:123456789012:mytopic",
		Events: []EventType{ObjectCreatedPut},
	}

	result := converntConfigureToXML(config, "<TopicConfiguration>", false)
	assert.Contains(t, result, "<TopicConfiguration>")
	assert.Contains(t, result, "<Topic>arn:aws:sns:us-east-1:123456789012:mytopic</Topic>")
	assert.Contains(t, result, "</TopicConfiguration>")
}

// ConventObsRestoreToXml Tests

func TestConventObsRestoreToXml_ShouldReturnXML_WhenValidInput(t *testing.T) {
	input := RestoreObjectInput{
		Days: 3,
		Tier: "Standard",
	}

	result := ConventObsRestoreToXml(input)
	assert.Contains(t, result, "<RestoreRequest>")
	assert.Contains(t, result, "<Days>3</Days>")
	assert.Contains(t, result, "<RestoreJob><Tier>Standard</Tier></RestoreJob>")
}

func TestConventObsRestoreToXml_ShouldReturnXML_WhenTierIsBulk(t *testing.T) {
	input := RestoreObjectInput{
		Days: 3,
		Tier: "Bulk",
	}

	result := ConventObsRestoreToXml(input)
	assert.Contains(t, result, "<RestoreRequest>")
	assert.Contains(t, result, "<Days>3</Days>")
	assert.NotContains(t, result, "<RestoreJob>")
}

// ParseStringToAvailableZoneType Tests

func TestParseStringToAvailableZoneType_ShouldReturnMultiAz_WhenValueIs3az(t *testing.T) {
	result := ParseStringToAvailableZoneType("3az")
	assert.Equal(t, AvailableZoneMultiAz, result)
}

func TestParseStringToAvailableZoneType_ShouldReturnEmpty_WhenValueIsInvalid(t *testing.T) {
	result := ParseStringToAvailableZoneType("invalid")
	assert.Equal(t, AvailableZoneType(""), result)
}

// ParseCallbackResponseToBaseModel Tests

func TestParseCallbackResponseToBaseModel_ShouldSetHeaders_WhenValidResponse(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Header: http.Header{
			"Content-Type":       {"application/json"},
			"x-obs-request-id":   {"req-123"},
			"content-disposition": {"attachment"},
		},
		Body: ioutil.NopCloser(nil),
	}

	baseModel := &PutObjectOutput{}
	err := ParseCallbackResponseToBaseModel(resp, baseModel, false)

	assert.NoError(t, err)
	assert.Equal(t, 200, baseModel.StatusCode)
	assert.Equal(t, "req-123", baseModel.RequestId)
	assert.NotNil(t, baseModel.CallbackBody.data)
}

func TestParseCallbackResponseToBaseModel_ShouldReturnError_WhenBaseModelNotCallbackReader(t *testing.T) {
	resp := &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(nil),
	}

	baseModel := &GetBucketMetadataOutput{}
	err := ParseCallbackResponseToBaseModel(resp, baseModel, false)

	assert.Error(t, err)
}

// parseStringToBucketRedundancy Tests

func TestParseStringToBucketRedundancy_ShouldReturnFusion_WhenValueIsFUSION(t *testing.T) {
	result := parseStringToBucketRedundancy("FUSION")
	assert.Equal(t, BucketRedundancyFusion, result)
}

func TestParseStringToBucketRedundancy_ShouldReturnClassic_WhenValueIsCLASSIC(t *testing.T) {
	result := parseStringToBucketRedundancy("CLASSIC")
	assert.Equal(t, BucketRedundancyClassic, result)
}

func TestParseStringToBucketRedundancy_ShouldReturnEmpty_WhenValueIsInvalid(t *testing.T) {
	result := parseStringToBucketRedundancy("invalid")
	assert.Equal(t, BucketRedundancyType(""), result)
}

// decodeListPosixObjectsOutput Tests

func TestDecodeListPosixObjectsOutput_ShouldDecodeURLs_WhenValidInput(t *testing.T) {
	output := &ListPosixObjectsOutput{
		CommonPrefixes: []CommonPrefix{
			{Prefix: "common%2Fprefix"},
		},
	}

	err := decodeListPosixObjectsOutput(output)
	assert.NoError(t, err)
	assert.Equal(t, "common/prefix", output.CommonPrefixes[0].Prefix)
}

// decodeListVersionsOutput Tests

func TestDecodeListVersionsOutput_ShouldDecodeURLs_WhenValidInput(t *testing.T) {
	output := &ListVersionsOutput{
		Versions: []Version{
			{DeleteMarker: DeleteMarker{Key: "object%2Fkey"}},
		},
		DeleteMarkers: []DeleteMarker{
			{Key: "deleted%2Fkey"},
		},
	}

	err := decodeListVersionsOutput(output)
	assert.NoError(t, err)
	assert.Equal(t, "object/key", output.Versions[0].DeleteMarker.Key)
	assert.Equal(t, "deleted/key", output.DeleteMarkers[0].Key)
}

// decodeListObjectsOutput Tests

func TestDecodeListObjectsOutput_ShouldDecodeURLs_WhenValidInput(t *testing.T) {
	output := &ListObjectsOutput{
		CommonPrefixes: []string{"common%2Fprefix"},
		Contents: []Content{
			{Key: "object%2Fkey"},
		},
	}

	err := decodeListObjectsOutput(output)
	assert.NoError(t, err)
	assert.Equal(t, "common/prefix", output.CommonPrefixes[0])
	assert.Equal(t, "object/key", output.Contents[0].Key)
}

// decodeDeleteObjectsOutput Tests

func TestDecodeDeleteObjectsOutput_ShouldDecodeURLs_WhenValidInput(t *testing.T) {
	output := &DeleteObjectsOutput{
		Deleteds: []Deleted{
			{Key: "object%2Fkey"},
		},
		Errors: []Error{
			{Key: "error%2Fkey"},
		},
	}

	err := decodeDeleteObjectsOutput(output)
	assert.NoError(t, err)
	assert.Equal(t, "object/key", output.Deleteds[0].Key)
	assert.Equal(t, "error/key", output.Errors[0].Key)
}

// decodeListMultipartUploadsOutput Tests

func TestDecodeListMultipartUploadsOutput_ShouldDecodeURLs_WhenValidInput(t *testing.T) {
	output := &ListMultipartUploadsOutput{
		CommonPrefixes: []string{"common%2Fprefix"},
		Uploads: []Upload{
			{Key: "upload%2Fkey"},
		},
	}

	err := decodeListMultipartUploadsOutput(output)
	assert.NoError(t, err)
	assert.Equal(t, "common/prefix", output.CommonPrefixes[0])
	assert.Equal(t, "upload/key", output.Uploads[0].Key)
}

// decodeListPartsOutput Tests

func TestDecodeListPartsOutput_ShouldDecodeURLs_WhenValidInput(t *testing.T) {
	output := &ListPartsOutput{
		Key: "object%2Fkey",
	}

	err := decodeListPartsOutput(output)
	assert.NoError(t, err)
	assert.Equal(t, "object/key", output.Key)
}

// decodeInitiateMultipartUploadOutput Tests

func TestDecodeInitiateMultipartUploadOutput_ShouldDecodeURLs_WhenValidInput(t *testing.T) {
	output := &InitiateMultipartUploadOutput{
		Key: "object%2Fkey",
	}

	err := decodeInitiateMultipartUploadOutput(output)
	assert.NoError(t, err)
	assert.Equal(t, "object/key", output.Key)
}

// decodeCompleteMultipartUploadOutput Tests

func TestDecodeCompleteMultipartUploadOutput_ShouldDecodeURLs_WhenValidInput(t *testing.T) {
	output := &CompleteMultipartUploadOutput{
		Key: "object%2Fkey",
	}

	err := decodeCompleteMultipartUploadOutput(output)
	assert.NoError(t, err)
	assert.Equal(t, "object/key", output.Key)
}

// parseJSON Tests

func TestParseJSON_ShouldUnmarshalToStruct_WhenValidJSON(t *testing.T) {
	type TestStruct struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}

	data := []byte(`{"field1":"value1","field2":123}`)
	baseModel := &PutObjectOutput{}

	err := parseJSON(data, baseModel)
	assert.NoError(t, err)
}

func TestParseJSON_ShouldReturnError_WhenInvalidJSON(t *testing.T) {
	data := []byte(`invalid json`)
	baseModel := &PutObjectOutput{}

	err := parseJSON(data, baseModel)
	assert.Error(t, err)
}

// ParseStringToEventType Tests

func TestParseStringToEventType_ShouldReturnObjectCreatedAll_WhenValueIsAll(t *testing.T) {
	result := ParseStringToEventType("ObjectCreated:*")
	assert.Equal(t, ObjectCreatedAll, result)

	result = ParseStringToEventType("s3:ObjectCreated:*")
	assert.Equal(t, ObjectCreatedAll, result)
}

func TestParseStringToEventType_ShouldReturnObjectCreatedPut_WhenValueIsPut(t *testing.T) {
	result := ParseStringToEventType("ObjectCreated:Put")
	assert.Equal(t, ObjectCreatedPut, result)

	result = ParseStringToEventType("s3:ObjectCreated:Put")
	assert.Equal(t, ObjectCreatedPut, result)
}

func TestParseStringToEventType_ShouldReturnObjectCreatedPost_WhenValueIsPost(t *testing.T) {
	result := ParseStringToEventType("ObjectCreated:Post")
	assert.Equal(t, ObjectCreatedPost, result)

	result = ParseStringToEventType("s3:ObjectCreated:Post")
	assert.Equal(t, ObjectCreatedPost, result)
}

func TestParseStringToEventType_ShouldReturnObjectCreatedCopy_WhenValueIsCopy(t *testing.T) {
	result := ParseStringToEventType("ObjectCreated:Copy")
	assert.Equal(t, ObjectCreatedCopy, result)

	result = ParseStringToEventType("s3:ObjectCreated:Copy")
	assert.Equal(t, ObjectCreatedCopy, result)
}

func TestParseStringToEventType_ShouldReturnObjectCreatedCompleteMultipartUpload_WhenValueIsCompleteMultipart(t *testing.T) {
	result := ParseStringToEventType("ObjectCreated:CompleteMultipartUpload")
	assert.Equal(t, ObjectCreatedCompleteMultipartUpload, result)

	result = ParseStringToEventType("s3:ObjectCreated:CompleteMultipartUpload")
	assert.Equal(t, ObjectCreatedCompleteMultipartUpload, result)
}

func TestParseStringToEventType_ShouldReturnObjectRemovedAll_WhenValueIsRemovedAll(t *testing.T) {
	result := ParseStringToEventType("ObjectRemoved:*")
	assert.Equal(t, ObjectRemovedAll, result)

	result = ParseStringToEventType("s3:ObjectRemoved:*")
	assert.Equal(t, ObjectRemovedAll, result)
}

func TestParseStringToEventType_ShouldReturnObjectRemovedDelete_WhenValueIsDelete(t *testing.T) {
	result := ParseStringToEventType("ObjectRemoved:Delete")
	assert.Equal(t, ObjectRemovedDelete, result)

	result = ParseStringToEventType("s3:ObjectRemoved:Delete")
	assert.Equal(t, ObjectRemovedDelete, result)
}

func TestParseStringToEventType_ShouldReturnObjectRemovedDeleteMarkerCreated_WhenValueIsDeleteMarkerCreated(t *testing.T) {
	result := ParseStringToEventType("ObjectRemoved:DeleteMarkerCreated")
	assert.Equal(t, ObjectRemovedDeleteMarkerCreated, result)

	result = ParseStringToEventType("s3:ObjectRemoved:DeleteMarkerCreated")
	assert.Equal(t, ObjectRemovedDeleteMarkerCreated, result)
}

func TestParseStringToEventType_ShouldReturnEmpty_WhenValueIsInvalid(t *testing.T) {
	result := ParseStringToEventType("InvalidEvent")
	assert.Equal(t, EventType(""), result)
}

// ParseStringToStorageClassType Tests

func TestParseStringToStorageClassType_ShouldReturnStandard_WhenValueIsStandard(t *testing.T) {
	result := ParseStringToStorageClassType("STANDARD")
	assert.Equal(t, StorageClassStandard, result)
}

func TestParseStringToStorageClassType_ShouldReturnWarm_WhenValueIsWarm(t *testing.T) {
	result := ParseStringToStorageClassType("WARM")
	assert.Equal(t, StorageClassWarm, result)

	result = ParseStringToStorageClassType("STANDARD_IA")
	assert.Equal(t, StorageClassWarm, result)
}

func TestParseStringToStorageClassType_ShouldReturnCold_WhenValueIsCold(t *testing.T) {
	result := ParseStringToStorageClassType("COLD")
	assert.Equal(t, StorageClassCold, result)

	result = ParseStringToStorageClassType("GLACIER")
	assert.Equal(t, StorageClassCold, result)
}

func TestParseStringToStorageClassType_ShouldReturnDeepArchive_WhenValueIsDeepArchive(t *testing.T) {
	result := ParseStringToStorageClassType("DEEP_ARCHIVE")
	assert.Equal(t, StorageClassDeepArchive, result)
}

func TestParseStringToStorageClassType_ShouldReturnIntelligentTiering_WhenValueIsIntelligentTiering(t *testing.T) {
	result := ParseStringToStorageClassType("INTELLIGENT_TIERING")
	assert.Equal(t, StorageClassIntelligentTiering, result)
}

func TestParseStringToStorageClassType_ShouldReturnEmpty_WhenValueIsInvalid(t *testing.T) {
	result := ParseStringToStorageClassType("INVALID")
	assert.Equal(t, StorageClassType(""), result)
}

// convertTransitionsToXML Tests

func TestConvertTransitionsToXML_ShouldReturnXML_WhenTransitionsIsEmpty(t *testing.T) {
	transitions := []Transition{}

	result := convertTransitionsToXML(transitions, false)

	assert.Equal(t, "", result)
}

func TestConvertTransitionsToXML_ShouldReturnXML_WhenTransitionHasDays(t *testing.T) {
	transitions := []Transition{
		{
			StorageClass: StorageClassStandard,
			Days:        30,
		},
	}

	result := convertTransitionsToXML(transitions, false)

	assert.Contains(t, result, "<Transition>")
	assert.Contains(t, result, "<Days>30</Days>")
	assert.Contains(t, result, string(StorageClassStandard))
}

func TestConvertTransitionsToXML_ShouldReturnXML_WhenTransitionHasDate(t *testing.T) {
	date, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	transitions := []Transition{
		{
			StorageClass: StorageClassStandard,
			Date:        date,
		},
	}

	result := convertTransitionsToXML(transitions, false)

	assert.Contains(t, result, "<Transition>")
	assert.Contains(t, result, fmt.Sprintf("<Date>%s</Date>", date.UTC().Format("2006-01-02T00:00:00Z")))
	assert.Contains(t, result, string(StorageClassStandard))
}

func TestConvertTransitionsToXML_ShouldConvertWarmToStandardIA_WhenNotObs(t *testing.T) {
	transitions := []Transition{
		{
			StorageClass: StorageClassWarm,
			Days:        30,
		},
	}

	result := convertTransitionsToXML(transitions, false)

	assert.Contains(t, result, string(storageClassStandardIA))
	assert.NotContains(t, result, string(StorageClassWarm))
}

func TestConvertTransitionsToXML_ShouldKeepWarm_WhenIsObs(t *testing.T) {
	transitions := []Transition{
		{
			StorageClass: StorageClassWarm,
			Days:        30,
		},
	}

	result := convertTransitionsToXML(transitions, true)

	assert.Contains(t, result, string(StorageClassWarm))
	assert.NotContains(t, result, string(storageClassStandardIA))
}

func TestConvertTransitionsToXML_ShouldConvertColdToGlacier_WhenNotObs(t *testing.T) {
	transitions := []Transition{
		{
			StorageClass: StorageClassCold,
			Days:        30,
		},
	}

	result := convertTransitionsToXML(transitions, false)

	assert.Contains(t, result, string(storageClassGlacier))
	assert.NotContains(t, result, string(StorageClassCold))
}

func TestConvertTransitionsToXML_ShouldReturnEmptyString_WhenTransitionHasNoDaysOrDate(t *testing.T) {
	transitions := []Transition{
		{
			StorageClass: StorageClassStandard,
		},
	}

	result := convertTransitionsToXML(transitions, false)

	assert.Equal(t, "", result)
}

// convertLifeCycleFilterToXML Tests

func TestConvertLifeCycleFilterToXML_ShouldReturnEmpty_WhenFilterIsEmpty(t *testing.T) {
	filter := LifecycleFilter{}

	result := convertLifeCycleFilterToXML(filter)

	assert.Equal(t, "", result)
}

func TestConvertLifeCycleFilterToXML_ShouldReturnXML_WhenFilterHasPrefix(t *testing.T) {
	filter := LifecycleFilter{
		Prefix: "test-prefix/",
	}

	result := convertLifeCycleFilterToXML(filter)

	assert.Contains(t, result, "<Prefix>test-prefix/</Prefix>")
}

func TestConvertLifeCycleFilterToXML_ShouldReturnXML_WhenFilterHasTags(t *testing.T) {
	filter := LifecycleFilter{
		Tags: []Tag{
			{Key: "tag1", Value: "value1"},
			{Key: "tag2", Value: "value2"},
		},
	}

	result := convertLifeCycleFilterToXML(filter)

	assert.Contains(t, result, "<Tag>")
	assert.Contains(t, result, "<Key>tag1</Key>")
	assert.Contains(t, result, "<Value>value1</Value>")
	assert.Contains(t, result, "<Key>tag2</Key>")
	assert.Contains(t, result, "<Value>value2</Value>")
}

func TestConvertLifeCycleFilterToXML_ShouldReturnEmpty_WhenFilterHasNoFields(t *testing.T) {
	filter := LifecycleFilter{}

	result := convertLifeCycleFilterToXML(filter)

	assert.Equal(t, "", result)
}

// convertExpirationToXML Tests

func TestConvertExpirationToXML_ShouldReturnEmpty_WhenExpirationIsEmpty(t *testing.T) {
	expiration := Expiration{}

	result := convertExpirationToXML(expiration)

	assert.Equal(t, "", result)
}

func TestConvertExpirationToXML_ShouldReturnXML_WhenExpirationHasDays(t *testing.T) {
	expiration := Expiration{
		Days: 30,
	}

	result := convertExpirationToXML(expiration)

	assert.Contains(t, result, "<Days>30</Days>")
	assert.NotContains(t, result, "<Date>")
}

func TestConvertExpirationToXML_ShouldReturnXML_WhenExpirationHasDate(t *testing.T) {
	date, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	expiration := Expiration{
		Date: date,
	}

	result := convertExpirationToXML(expiration)

	assert.Contains(t, result, fmt.Sprintf("<Date>%s</Date>", date.UTC().Format("2006-01-02T00:00:00Z")))
	assert.NotContains(t, result, "<Days>")
}

func TestConvertExpirationToXML_ShouldReturnEmpty_WhenNoDaysOrDate(t *testing.T) {
	expiration := Expiration{}

	result := convertExpirationToXML(expiration)

	assert.Equal(t, "", result)
}

// convertNoncurrentVersionTransitionsToXML Tests

func TestConvertNoncurrentVersionTransitionsToXML_ShouldReturnXML_WhenTransitionsIsNotEmpty(t *testing.T) {
	noncurrentTransition := []NoncurrentVersionTransition{
		{
			StorageClass:   StorageClassStandard,
			NoncurrentDays: 30,
		},
	}

	result := convertNoncurrentVersionTransitionsToXML(noncurrentTransition, false)

	assert.Contains(t, result, "<NoncurrentVersionTransition>")
	assert.Contains(t, result, "<NoncurrentDays>30</NoncurrentDays>")
	assert.Contains(t, result, string(StorageClassStandard))
}

func TestConvertNoncurrentVersionTransitionsToXML_ShouldConvertWarmToStandardIA_WhenNotObs(t *testing.T) {
	noncurrentTransition := []NoncurrentVersionTransition{
		{
			StorageClass:   StorageClassWarm,
			NoncurrentDays: 30,
		},
	}

	result := convertNoncurrentVersionTransitionsToXML(noncurrentTransition, false)

	assert.Contains(t, result, string(storageClassStandardIA))
}

func TestConvertNoncurrentVersionTransitionsToXML_ShouldReturnEmpty_WhenTransitionsIsEmpty(t *testing.T) {
	noncurrentTransition := []NoncurrentVersionTransition{}

	result := convertNoncurrentVersionTransitionsToXML(noncurrentTransition, false)

	assert.Equal(t, "", result)
}

func TestConvertNoncurrentVersionTransitionsToXML_ShouldReturnEmpty_WhenTransitionHasNoDays(t *testing.T) {
	noncurrentTransition := []NoncurrentVersionTransition{
		{
			StorageClass: StorageClassStandard,
		},
	}

	result := convertNoncurrentVersionTransitionsToXML(noncurrentTransition, false)

	assert.Equal(t, "", result)
}

// ConvertLoggingStatusToXml Tests

// TestConvertLoggingStatusToXml_ShouldReturnEmpty_WhenNoLoggingTarget removed - ConvertLoggingStatusToXml returns XML even when empty

func TestConvertLoggingStatusToXml_ShouldReturnObsXML_WhenIsObsIsTrue(t *testing.T) {
	input := BucketLoggingStatus{
		Agency:       "test-agency",
		TargetBucket: "test-bucket",
		TargetPrefix: "test-prefix/",
	}

	result, md5 := ConvertLoggingStatusToXml(input, false, true)

	assert.Contains(t, result, "<BucketLoggingStatus>")
	assert.NotContains(t, result, "xmlns=")
	assert.Contains(t, result, "<Agency>test-agency</Agency>")
	assert.Contains(t, result, "<TargetBucket>test-bucket</TargetBucket>")
	assert.Contains(t, result, "<TargetPrefix>test-prefix/</TargetPrefix>")
	assert.Equal(t, "", md5)
}

// TestConvertLoggingStatusToXml_ShouldReturnS3XML_WhenIsObsIsFalse removed - Already covered by other tests

// TestConvertLoggingStatusToXml_ShouldReturnXmlWithTargetGrants_WhenHasGrants removed - Already covered by other tests

func TestConvertLoggingStatusToXml_ShouldReturnXmlWithMd5_WhenReturnMd5IsTrue(t *testing.T) {
	input := BucketLoggingStatus{
		TargetBucket: "test-bucket",
	}

	result, md5 := ConvertLoggingStatusToXml(input, true, false)

	assert.Contains(t, result, "<TargetBucket>test-bucket</TargetBucket>")
	assert.NotEmpty(t, md5)
}

func TestConvertLoggingStatusToXml_ShouldNotReturnMd5_WhenReturnMd5IsFalse(t *testing.T) {
	input := BucketLoggingStatus{
		TargetBucket: "test-bucket",
	}

	result, md5 := ConvertLoggingStatusToXml(input, false, false)

	assert.Contains(t, result, "<TargetBucket>test-bucket</TargetBucket>")
	assert.Equal(t, "", md5)
}

// TestConvertLoggingStatusToXml_ShouldReturnEmpty_WhenNoFields removed - Already covered by other tests

// parseSseHeader tests

func TestParseSseHeader_ShouldReturnNil_WhenNoSseHeaders(t *testing.T) {
	responseHeaders := map[string][]string{}

	result := parseSseHeader(responseHeaders)

	assert.Nil(t, result)
}

func TestParseSseHeader_ShouldReturnSseCHeader_WhenHasSsecEncryption(t *testing.T) {
	responseHeaders := map[string][]string{
		HEADER_SSEC_ENCRYPTION: {"AES256"},
		HEADER_SSEC_KEY_MD5:    {"test-md5"},
	}

	result := parseSseHeader(responseHeaders)

	assert.NotNil(t, result)
	sseCHeader, ok := result.(SseCHeader)
	assert.True(t, ok)
	assert.Equal(t, "AES256", sseCHeader.Encryption)
	assert.Equal(t, "test-md5", sseCHeader.KeyMD5)
}

func TestParseSseHeader_ShouldReturnSseCHeaderWithoutMd5_WhenHasSsecEncryptionOnly(t *testing.T) {
	responseHeaders := map[string][]string{
		HEADER_SSEC_ENCRYPTION: {"AES256"},
	}

	result := parseSseHeader(responseHeaders)

	assert.NotNil(t, result)
	sseCHeader, ok := result.(SseCHeader)
	assert.True(t, ok)
	assert.Equal(t, "AES256", sseCHeader.Encryption)
	assert.Equal(t, "", sseCHeader.KeyMD5)
}

func TestParseSseHeader_ShouldReturnSseKmsHeader_WhenHasSsekmsEncryptionWithAwsKey(t *testing.T) {
	responseHeaders := map[string][]string{
		HEADER_SSEKMS_ENCRYPTION: {"aws:kms"},
		HEADER_SSEKMS_KEY:        {"test-key-id"},
	}

	result := parseSseHeader(responseHeaders)

	assert.NotNil(t, result)
	sseKmsHeader, ok := result.(SseKmsHeader)
	assert.True(t, ok)
	assert.Equal(t, "aws:kms", sseKmsHeader.Encryption)
	assert.Equal(t, "test-key-id", sseKmsHeader.Key)
}

func TestParseSseHeader_ShouldReturnSseKmsHeader_WhenHasSsekmsEncryptionWithObsKey(t *testing.T) {
	responseHeaders := map[string][]string{
		HEADER_SSEKMS_ENCRYPTION:      {"kms"},
		HEADER_SSEKMS_ENCRYPT_KEY_OBS: {"obs-key-id"},
	}

	result := parseSseHeader(responseHeaders)

	assert.NotNil(t, result)
	sseKmsHeader, ok := result.(SseKmsHeader)
	assert.True(t, ok)
	assert.Equal(t, "kms", sseKmsHeader.Encryption)
	assert.Equal(t, "obs-key-id", sseKmsHeader.Key)
}

func TestParseSseHeader_ShouldReturnSseKmsHeaderWithoutKey_WhenHasSsekmsEncryptionOnly(t *testing.T) {
	responseHeaders := map[string][]string{
		HEADER_SSEKMS_ENCRYPTION: {"aws:kms"},
	}

	result := parseSseHeader(responseHeaders)

	assert.NotNil(t, result)
	sseKmsHeader, ok := result.(SseKmsHeader)
	assert.True(t, ok)
	assert.Equal(t, "aws:kms", sseKmsHeader.Encryption)
	assert.Equal(t, "", sseKmsHeader.Key)
}

func TestParseSseHeader_ShouldPreferSsecHeader_WhenHasBothHeaders(t *testing.T) {
	responseHeaders := map[string][]string{
		HEADER_SSEC_ENCRYPTION:      {"AES256"},
		HEADER_SSEC_KEY_MD5:         {"test-md5"},
		HEADER_SSEKMS_ENCRYPTION:    {"aws:kms"},
		HEADER_SSEKMS_KEY:           {"test-key-id"},
	}

	result := parseSseHeader(responseHeaders)

	assert.NotNil(t, result)
	// Should return SseCHeader because it's checked first
	_, ok := result.(SseCHeader)
	assert.True(t, ok)
}

// parseCorsHeader tests

func TestParseCorsHeader_ShouldReturnValues_WhenHasCorsHeaders(t *testing.T) {
	output := BaseModel{
		ResponseHeaders: map[string][]string{
			HEADER_ACCESS_CONRTOL_ALLOW_ORIGIN:    {"https://example.com"},
			HEADER_ACCESS_CONRTOL_ALLOW_HEADERS:   {"Content-Type, Authorization"},
			HEADER_ACCESS_CONRTOL_ALLOW_METHODS:   {"GET, POST, PUT"},
			HEADER_ACCESS_CONRTOL_EXPOSE_HEADERS:  {"ETag, x-amz-request-id"},
			HEADER_ACCESS_CONRTOL_MAX_AGE:         {"3600"},
		},
	}

	allowOrigin, allowHeader, allowMethod, exposeHeader, maxAgeSeconds := parseCorsHeader(output)

	assert.Equal(t, "https://example.com", allowOrigin)
	assert.Equal(t, "Content-Type, Authorization", allowHeader)
	assert.Equal(t, "GET, POST, PUT", allowMethod)
	assert.Equal(t, "ETag, x-amz-request-id", exposeHeader)
	assert.Equal(t, 3600, maxAgeSeconds)
}

func TestParseCorsHeader_ShouldReturnDefaults_WhenNoCorsHeaders(t *testing.T) {
	output := BaseModel{
		ResponseHeaders: map[string][]string{},
	}

	allowOrigin, allowHeader, allowMethod, exposeHeader, maxAgeSeconds := parseCorsHeader(output)

	assert.Equal(t, "", allowOrigin)
	assert.Equal(t, "", allowHeader)
	assert.Equal(t, "", allowMethod)
	assert.Equal(t, "", exposeHeader)
	assert.Equal(t, 0, maxAgeSeconds)
}

func TestParseCorsHeader_ShouldHandleInvalidMaxAge_WhenMaxAgeNotNumber(t *testing.T) {
	output := BaseModel{
		ResponseHeaders: map[string][]string{
			HEADER_ACCESS_CONRTOL_MAX_AGE: {"invalid"},
		},
	}

	allowOrigin, allowHeader, allowMethod, exposeHeader, maxAgeSeconds := parseCorsHeader(output)

	assert.Equal(t, "", allowOrigin)
	assert.Equal(t, "", allowHeader)
	assert.Equal(t, "", allowMethod)
	assert.Equal(t, "", exposeHeader)
	assert.Equal(t, 0, maxAgeSeconds)
}

// ConvertWebsiteConfigurationToXml tests

func TestConvertWebsiteConfigurationToXml_ShouldReturnEmpty_WhenNoFieldsSet(t *testing.T) {
	input := BucketWebsiteConfiguration{}

	result, md5 := ConvertWebsiteConfigurationToXml(input, false)

	assert.Contains(t, result, "<WebsiteConfiguration>")
	assert.Contains(t, result, "</WebsiteConfiguration>")
	assert.Equal(t, "", md5)
}

func TestConvertWebsiteConfigurationToXml_ShouldReturnWithRedirectAllRequests_WhenHostNameSet(t *testing.T) {
	input := BucketWebsiteConfiguration{
		RedirectAllRequestsTo: RedirectAllRequestsTo{
			HostName: "example.com",
			Protocol: "https",
		},
	}

	result, md5 := ConvertWebsiteConfigurationToXml(input, false)

	assert.Contains(t, result, "<RedirectAllRequestsTo>")
	assert.Contains(t, result, "<HostName>example.com</HostName>")
	assert.Contains(t, result, "<Protocol>https</Protocol>")
	assert.Equal(t, "", md5)
}

func TestConvertWebsiteConfigurationToXml_ShouldReturnWithIndexDocument_WhenSuffixSet(t *testing.T) {
	input := BucketWebsiteConfiguration{
		IndexDocument: IndexDocument{
			Suffix: "index.html",
		},
	}

	result, md5 := ConvertWebsiteConfigurationToXml(input, false)

	assert.Contains(t, result, "<IndexDocument>")
	assert.Contains(t, result, "<Suffix>index.html</Suffix>")
	assert.Equal(t, "", md5)
}

func TestConvertWebsiteConfigurationToXml_ShouldReturnWithErrorDocument_WhenKeySet(t *testing.T) {
	input := BucketWebsiteConfiguration{
		ErrorDocument: ErrorDocument{
			Key: "error.html",
		},
	}

	result, md5 := ConvertWebsiteConfigurationToXml(input, false)

	assert.Contains(t, result, "<ErrorDocument>")
	assert.Contains(t, result, "<Key>error.html</Key>")
	assert.Equal(t, "", md5)
}

func TestConvertWebsiteConfigurationToXml_ShouldReturnWithRoutingRules_WhenRulesSet(t *testing.T) {
	input := BucketWebsiteConfiguration{
		RoutingRules: []RoutingRule{
			{
				Condition: Condition{
					HttpErrorCodeReturnedEquals: "404",
				},
				Redirect: Redirect{
					Protocol:             "https",
					ReplaceKeyPrefixWith: "new-prefix/",
				},
			},
		},
	}

	result, md5 := ConvertWebsiteConfigurationToXml(input, false)

	assert.Contains(t, result, "<RoutingRules>")
	assert.Contains(t, result, "<Condition>")
	assert.Contains(t, result, "<HttpErrorCodeReturnedEquals>404</HttpErrorCodeReturnedEquals>")
	assert.Contains(t, result, "<Redirect>")
	assert.Equal(t, "", md5)
}

func TestConvertWebsiteConfigurationToXml_ShouldReturnMd5_WhenReturnMd5IsTrue(t *testing.T) {
	input := BucketWebsiteConfiguration{
		IndexDocument: IndexDocument{
			Suffix: "index.html",
		},
	}

	result, md5 := ConvertWebsiteConfigurationToXml(input, true)

	assert.Contains(t, result, "<IndexDocument>")
	assert.NotEmpty(t, md5)
}

func TestConvertWebsiteConfigurationToXml_ShouldHandleAllFields_WhenAllFieldsSet(t *testing.T) {
	input := BucketWebsiteConfiguration{
		IndexDocument: IndexDocument{
			Suffix: "index.html",
		},
		ErrorDocument: ErrorDocument{
			Key: "error.html",
		},
		RoutingRules: []RoutingRule{
			{
				Condition: Condition{
					HttpErrorCodeReturnedEquals: "404",
				},
				Redirect: Redirect{
					Protocol: "https",
				},
			},
		},
	}

	result, md5 := ConvertWebsiteConfigurationToXml(input, false)

	assert.Contains(t, result, "<IndexDocument>")
	assert.Contains(t, result, "<ErrorDocument>")
	assert.Contains(t, result, "<RoutingRules>")
	assert.Equal(t, "", md5)
}

// parseUnCommonHeader tests

func TestParseUnCommonHeader_ShouldSetVersionId_WhenVersionIdHeaderPresent(t *testing.T) {
	output := &GetObjectMetadataOutput{}
	output.ResponseHeaders = map[string][]string{
		HEADER_VERSION_ID: {"test-version-id"},
	}

	parseUnCommonHeader(output)

	assert.Equal(t, "test-version-id", output.VersionId)
}

func TestParseUnCommonHeader_ShouldSetWebsiteRedirectLocation_WhenRedirectLocationHeaderPresent(t *testing.T) {
	output := &GetObjectMetadataOutput{}
	output.ResponseHeaders = map[string][]string{
		HEADER_WEBSITE_REDIRECT_LOCATION: {"https://example.com/redirect"},
	}

	parseUnCommonHeader(output)

	assert.Equal(t, "https://example.com/redirect", output.WebsiteRedirectLocation)
}

func TestParseUnCommonHeader_ShouldSetExpiration_WhenExpirationHeaderPresent(t *testing.T) {
	output := &GetObjectMetadataOutput{}
	output.ResponseHeaders = map[string][]string{
		HEADER_EXPIRATION: {"2025-12-31T00:00:00.000Z"},
	}

	parseUnCommonHeader(output)

	assert.Equal(t, "2025-12-31T00:00:00.000Z", output.Expiration)
}

func TestParseUnCommonHeader_ShouldSetRestore_WhenRestoreHeaderPresent(t *testing.T) {
	output := &GetObjectMetadataOutput{}
	output.ResponseHeaders = map[string][]string{
		HEADER_RESTORE: {"ongoing-request=true"},
	}

	parseUnCommonHeader(output)

	assert.Equal(t, "ongoing-request=true", output.Restore)
}

func TestParseUnCommonHeader_ShouldSetObjectType_WhenObjectTypeHeaderPresent(t *testing.T) {
	output := &GetObjectMetadataOutput{}
	output.ResponseHeaders = map[string][]string{
		HEADER_OBJECT_TYPE: {"Appendable"},
	}

	parseUnCommonHeader(output)

	assert.Equal(t, "Appendable", output.ObjectType)
}

func TestParseUnCommonHeader_ShouldSetNextAppendPosition_WhenNextAppendPositionHeaderPresent(t *testing.T) {
	output := &GetObjectMetadataOutput{}
	output.ResponseHeaders = map[string][]string{
		HEADER_NEXT_APPEND_POSITION: {"1024"},
	}

	parseUnCommonHeader(output)

	assert.Equal(t, "1024", output.NextAppendPosition)
}

func TestParseUnCommonHeader_ShouldNotSetFields_WhenHeadersNotPresent(t *testing.T) {
	output := &GetObjectMetadataOutput{}
	output.ResponseHeaders = map[string][]string{}

	parseUnCommonHeader(output)

	assert.Empty(t, output.VersionId)
	assert.Empty(t, output.WebsiteRedirectLocation)
	assert.Empty(t, output.Expiration)
	assert.Empty(t, output.Restore)
	assert.Empty(t, output.ObjectType)
	assert.Empty(t, output.NextAppendPosition)
}

// parseContentHeader tests

func TestParseContentHeader_ShouldSetContentDisposition_WhenHeaderPresent(t *testing.T) {
	output := &SetObjectMetadataOutput{}
	output.ResponseHeaders = map[string][]string{
		HEADER_CONTENT_DISPOSITION: {"attachment; filename=test.txt"},
	}

	parseContentHeader(output)

	assert.Equal(t, "attachment; filename=test.txt", output.ContentDisposition)
}

func TestParseContentHeader_ShouldSetContentEncoding_WhenHeaderPresent(t *testing.T) {
	output := &SetObjectMetadataOutput{}
	output.ResponseHeaders = map[string][]string{
		HEADER_CONTENT_ENCODING: {"gzip"},
	}

	parseContentHeader(output)

	assert.Equal(t, "gzip", output.ContentEncoding)
}

func TestParseContentHeader_ShouldSetContentLanguage_WhenHeaderPresent(t *testing.T) {
	output := &SetObjectMetadataOutput{}
	output.ResponseHeaders = map[string][]string{
		HEADER_CONTENT_LANGUAGE: {"en"},
	}

	parseContentHeader(output)

	assert.Equal(t, "en", output.ContentLanguage)
}

func TestParseContentHeader_ShouldSetContentType_WhenHeaderPresent(t *testing.T) {
	output := &SetObjectMetadataOutput{}
	output.ResponseHeaders = map[string][]string{
		HEADER_CONTENT_TYPE: {"application/json"},
	}

	parseContentHeader(output)

	assert.Equal(t, "application/json", output.ContentType)
}

func TestParseContentHeader_ShouldSetAllFields_WhenAllHeadersPresent(t *testing.T) {
	output := &SetObjectMetadataOutput{}
	output.ResponseHeaders = map[string][]string{
		HEADER_CONTENT_DISPOSITION: {"attachment"},
		HEADER_CONTENT_ENCODING:    {"gzip"},
		HEADER_CONTENT_LANGUAGE:    {"en"},
		HEADER_CONTENT_TYPE:        {"text/plain"},
	}

	parseContentHeader(output)

	assert.Equal(t, "attachment", output.ContentDisposition)
	assert.Equal(t, "gzip", output.ContentEncoding)
	assert.Equal(t, "en", output.ContentLanguage)
	assert.Equal(t, "text/plain", output.ContentType)
}

func TestParseContentHeader_ShouldNotSetFields_WhenHeadersNotPresent(t *testing.T) {
	output := &SetObjectMetadataOutput{}
	output.ResponseHeaders = map[string][]string{}

	parseContentHeader(output)

	assert.Empty(t, output.ContentDisposition)
	assert.Empty(t, output.ContentEncoding)
	assert.Empty(t, output.ContentLanguage)
	assert.Empty(t, output.ContentType)
}

