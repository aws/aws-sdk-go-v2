package sqs

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go/middleware"
)

//********************************
// SendMessage checksum validation
//********************************
func addValidateSendMessageChecksum(stack *middleware.Stack, o Options) error {
	return addValidateMessageChecksum(stack, o, validateSendMessageChecksum)
}

func validateSendMessageChecksum(input, output interface{}) error {
	in, ok := input.(*SendMessageInput)
	if !ok {
		return fmt.Errorf("wrong input type, expect %T, got %T", in, input)
	}
	out, ok := output.(*SendMessageOutput)
	if !ok {
		return fmt.Errorf("wrong output type, expect %T, got %T", out, output)
	}

	// Nothing to validate if the members aren't populated.
	if in.MessageBody == nil || out.MD5OfMessageBody == nil {
		return nil
	}

	if err := validateMessageChecksum(*in.MessageBody, *out.MD5OfMessageBody); err != nil {
		return messageChecksumError{
			MessageID: aws.ToString(out.MessageId),
			Err:       err,
		}
	}
	return nil
}

//*************************************
// SendMessageBatch checksum validation
//*************************************
func addValidateSendMessageBatchChecksum(stack *middleware.Stack, o Options) error {
	return addValidateMessageChecksum(stack, o, validateSendMessageBatchChecksum)
}

func validateSendMessageBatchChecksum(input, output interface{}) error {
	in, ok := input.(*SendMessageBatchInput)
	if !ok {
		return fmt.Errorf("wrong input type, expect %T, got %T", in, input)
	}
	out, ok := output.(*SendMessageBatchOutput)
	if !ok {
		return fmt.Errorf("wrong output type, expect %T, got %T", out, output)
	}

	outEntries := map[string]sqstypes.SendMessageBatchResultEntry{}
	for _, e := range out.Successful {
		outEntries[*e.Id] = e
	}

	var failedMessageErrs []error
	for _, inEntry := range in.Entries {
		outEntry, ok := outEntries[*inEntry.Id]
		// Nothing to validate if the members aren't populated.
		if !ok || inEntry.MessageBody == nil || outEntry.MD5OfMessageBody == nil {
			continue
		}

		if err := validateMessageChecksum(*inEntry.MessageBody, *outEntry.MD5OfMessageBody); err != nil {
			failedMessageErrs = append(failedMessageErrs, messageChecksumError{
				MessageID: aws.ToString(outEntry.MessageId),
				Err:       err,
			})
		}
	}

	if len(failedMessageErrs) != 0 {
		return batchMessageChecksumError{
			Errs: failedMessageErrs,
		}
	}

	return nil
}

//***********************************
// ReceiveMessage checksum validation
//***********************************
func addValidateReceiveMessageChecksum(stack *middleware.Stack, o Options) error {
	return addValidateMessageChecksum(stack, o, validateReceiveMessageChecksum)
}

func validateReceiveMessageChecksum(_, output interface{}) error {
	out, ok := output.(*ReceiveMessageOutput)
	if !ok {
		return fmt.Errorf("wrong output type, expect %T, got %T", out, output)
	}

	var failedMessageErrs []error
	for _, msg := range out.Messages {
		// Nothing to validate if the members aren't populated.
		if msg.Body == nil || msg.MD5OfBody == nil {
			continue
		}

		if err := validateMessageChecksum(*msg.Body, *msg.MD5OfBody); err != nil {
			failedMessageErrs = append(failedMessageErrs, messageChecksumError{
				MessageID: aws.ToString(msg.MessageId),
				Err:       err,
			})
		}
	}

	if len(failedMessageErrs) != 0 {
		return batchMessageChecksumError{
			Errs: failedMessageErrs,
		}
	}

	return nil
}

//***************************************
// Message checksum validation middleware
//***************************************
type messageChecksumValidator func(input, output interface{}) error

func addValidateMessageChecksum(stack *middleware.Stack, o Options, validate messageChecksumValidator) error {
	if o.DisableMessageChecksumValidation {
		return nil
	}

	m := validateMessageChecksumMiddleware{
		validate: validate,
	}
	err := stack.Initialize.Add(m, middleware.Before)
	if err != nil {
		return fmt.Errorf("failed to add %s middleware, %w", m.ID(), err)
	}

	return nil
}

type validateMessageChecksumMiddleware struct {
	validate messageChecksumValidator
}

func (validateMessageChecksumMiddleware) ID() string { return "SQSValidateMessageChecksum" }

func (m validateMessageChecksumMiddleware) HandleInitialize(
	ctx context.Context, input middleware.InitializeInput, next middleware.InitializeHandler,
) (
	out middleware.InitializeOutput, meta middleware.Metadata, err error,
) {
	out, meta, err = next.HandleInitialize(ctx, input)
	if err != nil {
		return out, meta, err
	}

	err = m.validate(input.Parameters, out.Result)
	if err != nil {
		return out, meta, fmt.Errorf("message checksum validation failed, %w", err)
	}

	return out, meta, nil
}

func validateMessageChecksum(value, expect string) error {
	msum := md5.Sum([]byte(value))
	sum := hex.EncodeToString(msum[:])
	if sum != expect {
		return fmt.Errorf("expected MD5 checksum %s, got %s", expect, sum)
	}

	return nil
}

//************************
// Message checksum errors
//************************
type messageChecksumError struct {
	MessageID string
	Err       error
}

func (e messageChecksumError) Error() string {
	prefix := "message"
	if e.MessageID != "" {
		prefix += " " + e.MessageID
	}
	return fmt.Sprintf("%s has invalid checksum, %v", prefix, e.Err.Error())
}

type batchMessageChecksumError struct {
	Errs []error
}

func (e batchMessageChecksumError) Error() string {
	var w strings.Builder
	fmt.Fprintf(&w, "message checksum errors")

	for _, err := range e.Errs {
		fmt.Fprintf(&w, "\n\t%s", err.Error())
	}

	return w.String()
}
