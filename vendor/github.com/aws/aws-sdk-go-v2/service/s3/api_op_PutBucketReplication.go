// Code generated by smithy-go-codegen DO NOT EDIT.

package s3

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	internalauth "github.com/aws/aws-sdk-go-v2/internal/auth"
	"github.com/aws/aws-sdk-go-v2/internal/v4a"
	internalChecksum "github.com/aws/aws-sdk-go-v2/service/internal/checksum"
	s3cust "github.com/aws/aws-sdk-go-v2/service/s3/internal/customizations"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// Creates a replication configuration or replaces an existing one. For more
// information, see Replication (https://docs.aws.amazon.com/AmazonS3/latest/dev/replication.html)
// in the Amazon S3 User Guide. Specify the replication configuration in the
// request body. In the replication configuration, you provide the name of the
// destination bucket or buckets where you want Amazon S3 to replicate objects, the
// IAM role that Amazon S3 can assume to replicate objects on your behalf, and
// other relevant information. A replication configuration must include at least
// one rule, and can contain a maximum of 1,000. Each rule identifies a subset of
// objects to replicate by filtering the objects in the source bucket. To choose
// additional subsets of objects to replicate, add a rule for each subset. To
// specify a subset of the objects in the source bucket to apply a replication rule
// to, add the Filter element as a child of the Rule element. You can filter
// objects based on an object key prefix, one or more object tags, or both. When
// you add the Filter element in the configuration, you must also add the following
// elements: DeleteMarkerReplication , Status , and Priority . If you are using an
// earlier version of the replication configuration, Amazon S3 handles replication
// of delete markers differently. For more information, see Backward Compatibility (https://docs.aws.amazon.com/AmazonS3/latest/dev/replication-add-config.html#replication-backward-compat-considerations)
// . For information about enabling versioning on a bucket, see Using Versioning (https://docs.aws.amazon.com/AmazonS3/latest/dev/Versioning.html)
// . Handling Replication of Encrypted Objects By default, Amazon S3 doesn't
// replicate objects that are stored at rest using server-side encryption with KMS
// keys. To replicate Amazon Web Services KMS-encrypted objects, add the following:
// SourceSelectionCriteria , SseKmsEncryptedObjects , Status ,
// EncryptionConfiguration , and ReplicaKmsKeyID . For information about
// replication configuration, see Replicating Objects Created with SSE Using KMS
// keys (https://docs.aws.amazon.com/AmazonS3/latest/dev/replication-config-for-kms-objects.html)
// . For information on PutBucketReplication errors, see List of
// replication-related error codes (https://docs.aws.amazon.com/AmazonS3/latest/API/ErrorResponses.html#ReplicationErrorCodeList)
// Permissions To create a PutBucketReplication request, you must have
// s3:PutReplicationConfiguration permissions for the bucket. By default, a
// resource owner, in this case the Amazon Web Services account that created the
// bucket, can perform this operation. The resource owner can also grant others
// permissions to perform the operation. For more information about permissions,
// see Specifying Permissions in a Policy (https://docs.aws.amazon.com/AmazonS3/latest/dev/using-with-s3-actions.html)
// and Managing Access Permissions to Your Amazon S3 Resources (https://docs.aws.amazon.com/AmazonS3/latest/userguide/s3-access-control.html)
// . To perform this operation, the user or role performing the action must have
// the iam:PassRole (https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_use_passrole.html)
// permission. The following operations are related to PutBucketReplication :
//   - GetBucketReplication (https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketReplication.html)
//   - DeleteBucketReplication (https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketReplication.html)
func (c *Client) PutBucketReplication(ctx context.Context, params *PutBucketReplicationInput, optFns ...func(*Options)) (*PutBucketReplicationOutput, error) {
	if params == nil {
		params = &PutBucketReplicationInput{}
	}

	result, metadata, err := c.invokeOperation(ctx, "PutBucketReplication", params, optFns, c.addOperationPutBucketReplicationMiddlewares)
	if err != nil {
		return nil, err
	}

	out := result.(*PutBucketReplicationOutput)
	out.ResultMetadata = metadata
	return out, nil
}

type PutBucketReplicationInput struct {

	// The name of the bucket
	//
	// This member is required.
	Bucket *string

	// A container for replication rules. You can add up to 1,000 rules. The maximum
	// size of a replication configuration is 2 MB.
	//
	// This member is required.
	ReplicationConfiguration *types.ReplicationConfiguration

	// Indicates the algorithm used to create the checksum for the object when using
	// the SDK. This header will not provide any additional functionality if not using
	// the SDK. When sending this header, there must be a corresponding x-amz-checksum
	// or x-amz-trailer header sent. Otherwise, Amazon S3 fails the request with the
	// HTTP status code 400 Bad Request . For more information, see Checking object
	// integrity (https://docs.aws.amazon.com/AmazonS3/latest/userguide/checking-object-integrity.html)
	// in the Amazon S3 User Guide. If you provide an individual checksum, Amazon S3
	// ignores any provided ChecksumAlgorithm parameter.
	ChecksumAlgorithm types.ChecksumAlgorithm

	// The base64-encoded 128-bit MD5 digest of the data. You must use this header as
	// a message integrity check to verify that the request body was not corrupted in
	// transit. For more information, see RFC 1864 (http://www.ietf.org/rfc/rfc1864.txt)
	// . For requests made using the Amazon Web Services Command Line Interface (CLI)
	// or Amazon Web Services SDKs, this field is calculated automatically.
	ContentMD5 *string

	// The account ID of the expected bucket owner. If the bucket is owned by a
	// different account, the request fails with the HTTP status code 403 Forbidden
	// (access denied).
	ExpectedBucketOwner *string

	// A token to allow Object Lock to be enabled for an existing bucket.
	Token *string

	noSmithyDocumentSerde
}

type PutBucketReplicationOutput struct {
	// Metadata pertaining to the operation's result.
	ResultMetadata middleware.Metadata

	noSmithyDocumentSerde
}

func (c *Client) addOperationPutBucketReplicationMiddlewares(stack *middleware.Stack, options Options) (err error) {
	err = stack.Serialize.Add(&awsRestxml_serializeOpPutBucketReplication{}, middleware.After)
	if err != nil {
		return err
	}
	err = stack.Deserialize.Add(&awsRestxml_deserializeOpPutBucketReplication{}, middleware.After)
	if err != nil {
		return err
	}
	if err = addlegacyEndpointContextSetter(stack, options); err != nil {
		return err
	}
	if err = addSetLoggerMiddleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddClientRequestIDMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddComputeContentLengthMiddleware(stack); err != nil {
		return err
	}
	if err = addResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = v4.AddComputePayloadSHA256Middleware(stack); err != nil {
		return err
	}
	if err = addRetryMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addHTTPSignerV4Middleware(stack, options); err != nil {
		return err
	}
	if err = awsmiddleware.AddRawResponseToMetadata(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecordResponseTiming(stack); err != nil {
		return err
	}
	if err = addClientUserAgent(stack, options); err != nil {
		return err
	}
	if err = smithyhttp.AddErrorCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = smithyhttp.AddCloseResponseBodyMiddleware(stack); err != nil {
		return err
	}
	if err = swapWithCustomHTTPSignerMiddleware(stack, options); err != nil {
		return err
	}
	if err = addPutBucketReplicationResolveEndpointMiddleware(stack, options); err != nil {
		return err
	}
	if err = addOpPutBucketReplicationValidationMiddleware(stack); err != nil {
		return err
	}
	if err = stack.Initialize.Add(newServiceMetadataMiddleware_opPutBucketReplication(options.Region), middleware.Before); err != nil {
		return err
	}
	if err = addMetadataRetrieverMiddleware(stack); err != nil {
		return err
	}
	if err = awsmiddleware.AddRecursionDetection(stack); err != nil {
		return err
	}
	if err = addPutBucketReplicationInputChecksumMiddlewares(stack, options); err != nil {
		return err
	}
	if err = addPutBucketReplicationUpdateEndpoint(stack, options); err != nil {
		return err
	}
	if err = addResponseErrorMiddleware(stack); err != nil {
		return err
	}
	if err = v4.AddContentSHA256HeaderMiddleware(stack); err != nil {
		return err
	}
	if err = disableAcceptEncodingGzip(stack); err != nil {
		return err
	}
	if err = addRequestResponseLogging(stack, options); err != nil {
		return err
	}
	if err = addendpointDisableHTTPSMiddleware(stack, options); err != nil {
		return err
	}
	if err = addSerializeImmutableHostnameBucketMiddleware(stack, options); err != nil {
		return err
	}
	return nil
}

func (v *PutBucketReplicationInput) bucket() (string, bool) {
	if v.Bucket == nil {
		return "", false
	}
	return *v.Bucket, true
}

func newServiceMetadataMiddleware_opPutBucketReplication(region string) *awsmiddleware.RegisterServiceMetadata {
	return &awsmiddleware.RegisterServiceMetadata{
		Region:        region,
		ServiceID:     ServiceID,
		SigningName:   "s3",
		OperationName: "PutBucketReplication",
	}
}

// getPutBucketReplicationRequestAlgorithmMember gets the request checksum
// algorithm value provided as input.
func getPutBucketReplicationRequestAlgorithmMember(input interface{}) (string, bool) {
	in := input.(*PutBucketReplicationInput)
	if len(in.ChecksumAlgorithm) == 0 {
		return "", false
	}
	return string(in.ChecksumAlgorithm), true
}

func addPutBucketReplicationInputChecksumMiddlewares(stack *middleware.Stack, options Options) error {
	return internalChecksum.AddInputMiddleware(stack, internalChecksum.InputMiddlewareOptions{
		GetAlgorithm:                     getPutBucketReplicationRequestAlgorithmMember,
		RequireChecksum:                  true,
		EnableTrailingChecksum:           false,
		EnableComputeSHA256PayloadHash:   true,
		EnableDecodedContentLengthHeader: true,
	})
}

// getPutBucketReplicationBucketMember returns a pointer to string denoting a
// provided bucket member valueand a boolean indicating if the input has a modeled
// bucket name,
func getPutBucketReplicationBucketMember(input interface{}) (*string, bool) {
	in := input.(*PutBucketReplicationInput)
	if in.Bucket == nil {
		return nil, false
	}
	return in.Bucket, true
}
func addPutBucketReplicationUpdateEndpoint(stack *middleware.Stack, options Options) error {
	return s3cust.UpdateEndpoint(stack, s3cust.UpdateEndpointOptions{
		Accessor: s3cust.UpdateEndpointParameterAccessor{
			GetBucketFromInput: getPutBucketReplicationBucketMember,
		},
		UsePathStyle:                   options.UsePathStyle,
		UseAccelerate:                  options.UseAccelerate,
		SupportsAccelerate:             true,
		TargetS3ObjectLambda:           false,
		EndpointResolver:               options.EndpointResolver,
		EndpointResolverOptions:        options.EndpointOptions,
		UseARNRegion:                   options.UseARNRegion,
		DisableMultiRegionAccessPoints: options.DisableMultiRegionAccessPoints,
	})
}

type opPutBucketReplicationResolveEndpointMiddleware struct {
	EndpointResolver EndpointResolverV2
	BuiltInResolver  builtInParameterResolver
}

func (*opPutBucketReplicationResolveEndpointMiddleware) ID() string {
	return "ResolveEndpointV2"
}

func (m *opPutBucketReplicationResolveEndpointMiddleware) HandleSerialize(ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {
	if awsmiddleware.GetRequiresLegacyEndpoints(ctx) {
		return next.HandleSerialize(ctx, in)
	}

	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", in.Request)
	}

	input, ok := in.Parameters.(*PutBucketReplicationInput)
	if !ok {
		return out, metadata, fmt.Errorf("unknown transport type %T", in.Request)
	}

	if m.EndpointResolver == nil {
		return out, metadata, fmt.Errorf("expected endpoint resolver to not be nil")
	}

	params := EndpointParameters{}

	m.BuiltInResolver.ResolveBuiltIns(&params)

	params.Bucket = input.Bucket

	var resolvedEndpoint smithyendpoints.Endpoint
	resolvedEndpoint, err = m.EndpointResolver.ResolveEndpoint(ctx, params)
	if err != nil {
		return out, metadata, fmt.Errorf("failed to resolve service endpoint, %w", err)
	}

	req.URL = &resolvedEndpoint.URI

	for k := range resolvedEndpoint.Headers {
		req.Header.Set(
			k,
			resolvedEndpoint.Headers.Get(k),
		)
	}

	authSchemes, err := internalauth.GetAuthenticationSchemes(&resolvedEndpoint.Properties)
	if err != nil {
		var nfe *internalauth.NoAuthenticationSchemesFoundError
		if errors.As(err, &nfe) {
			// if no auth scheme is found, default to sigv4
			signingName := "s3"
			signingRegion := m.BuiltInResolver.(*builtInResolver).Region
			ctx = awsmiddleware.SetSigningName(ctx, signingName)
			ctx = awsmiddleware.SetSigningRegion(ctx, signingRegion)
			ctx = s3cust.SetSignerVersion(ctx, internalauth.SigV4)
		}
		var ue *internalauth.UnSupportedAuthenticationSchemeSpecifiedError
		if errors.As(err, &ue) {
			return out, metadata, fmt.Errorf(
				"This operation requests signer version(s) %v but the client only supports %v",
				ue.UnsupportedSchemes,
				internalauth.SupportedSchemes,
			)
		}
	}

	for _, authScheme := range authSchemes {
		switch authScheme.(type) {
		case *internalauth.AuthenticationSchemeV4:
			v4Scheme, _ := authScheme.(*internalauth.AuthenticationSchemeV4)
			var signingName, signingRegion string
			if v4Scheme.SigningName == nil {
				signingName = "s3"
			} else {
				signingName = *v4Scheme.SigningName
			}
			if v4Scheme.SigningRegion == nil {
				signingRegion = m.BuiltInResolver.(*builtInResolver).Region
			} else {
				signingRegion = *v4Scheme.SigningRegion
			}
			if v4Scheme.DisableDoubleEncoding != nil {
				// The signer sets an equivalent value at client initialization time.
				// Setting this context value will cause the signer to extract it
				// and override the value set at client initialization time.
				ctx = internalauth.SetDisableDoubleEncoding(ctx, *v4Scheme.DisableDoubleEncoding)
			}
			ctx = awsmiddleware.SetSigningName(ctx, signingName)
			ctx = awsmiddleware.SetSigningRegion(ctx, signingRegion)
			ctx = s3cust.SetSignerVersion(ctx, v4Scheme.Name)
			break
		case *internalauth.AuthenticationSchemeV4A:
			v4aScheme, _ := authScheme.(*internalauth.AuthenticationSchemeV4A)
			if v4aScheme.SigningName == nil {
				v4aScheme.SigningName = aws.String("s3")
			}
			if v4aScheme.DisableDoubleEncoding != nil {
				// The signer sets an equivalent value at client initialization time.
				// Setting this context value will cause the signer to extract it
				// and override the value set at client initialization time.
				ctx = internalauth.SetDisableDoubleEncoding(ctx, *v4aScheme.DisableDoubleEncoding)
			}
			ctx = awsmiddleware.SetSigningName(ctx, *v4aScheme.SigningName)
			ctx = awsmiddleware.SetSigningRegion(ctx, v4aScheme.SigningRegionSet[0])
			ctx = s3cust.SetSignerVersion(ctx, v4a.Version)
			break
		case *internalauth.AuthenticationSchemeNone:
			break
		}
	}

	return next.HandleSerialize(ctx, in)
}

func addPutBucketReplicationResolveEndpointMiddleware(stack *middleware.Stack, options Options) error {
	return stack.Serialize.Insert(&opPutBucketReplicationResolveEndpointMiddleware{
		EndpointResolver: options.EndpointResolverV2,
		BuiltInResolver: &builtInResolver{
			Region:                         options.Region,
			UseFIPS:                        options.EndpointOptions.UseFIPSEndpoint,
			UseDualStack:                   options.EndpointOptions.UseDualStackEndpoint,
			Endpoint:                       options.BaseEndpoint,
			ForcePathStyle:                 options.UsePathStyle,
			Accelerate:                     options.UseAccelerate,
			DisableMultiRegionAccessPoints: options.DisableMultiRegionAccessPoints,
			UseArnRegion:                   options.UseARNRegion,
		},
	}, "ResolveEndpoint", middleware.After)
}
