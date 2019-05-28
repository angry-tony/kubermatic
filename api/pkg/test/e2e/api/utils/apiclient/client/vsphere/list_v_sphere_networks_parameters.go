// Code generated by go-swagger; DO NOT EDIT.

package vsphere

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"
)

// NewListVSphereNetworksParams creates a new ListVSphereNetworksParams object
// with the default values initialized.
func NewListVSphereNetworksParams() *ListVSphereNetworksParams {

	return &ListVSphereNetworksParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewListVSphereNetworksParamsWithTimeout creates a new ListVSphereNetworksParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewListVSphereNetworksParamsWithTimeout(timeout time.Duration) *ListVSphereNetworksParams {

	return &ListVSphereNetworksParams{

		timeout: timeout,
	}
}

// NewListVSphereNetworksParamsWithContext creates a new ListVSphereNetworksParams object
// with the default values initialized, and the ability to set a context for a request
func NewListVSphereNetworksParamsWithContext(ctx context.Context) *ListVSphereNetworksParams {

	return &ListVSphereNetworksParams{

		Context: ctx,
	}
}

// NewListVSphereNetworksParamsWithHTTPClient creates a new ListVSphereNetworksParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewListVSphereNetworksParamsWithHTTPClient(client *http.Client) *ListVSphereNetworksParams {

	return &ListVSphereNetworksParams{
		HTTPClient: client,
	}
}

/*ListVSphereNetworksParams contains all the parameters to send to the API endpoint
for the list v sphere networks operation typically these are written to a http.Request
*/
type ListVSphereNetworksParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the list v sphere networks params
func (o *ListVSphereNetworksParams) WithTimeout(timeout time.Duration) *ListVSphereNetworksParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list v sphere networks params
func (o *ListVSphereNetworksParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list v sphere networks params
func (o *ListVSphereNetworksParams) WithContext(ctx context.Context) *ListVSphereNetworksParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list v sphere networks params
func (o *ListVSphereNetworksParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list v sphere networks params
func (o *ListVSphereNetworksParams) WithHTTPClient(client *http.Client) *ListVSphereNetworksParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list v sphere networks params
func (o *ListVSphereNetworksParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *ListVSphereNetworksParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}