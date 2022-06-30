# \EventServiceApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**EventServiceCreateEvent**](EventServiceApi.md#EventServiceCreateEvent) | **Post** /event.EventService/CreateEvent | 
[**EventServiceDeleteEvent**](EventServiceApi.md#EventServiceDeleteEvent) | **Post** /event.EventService/DeleteEvent | 
[**EventServiceGetEventsForDay**](EventServiceApi.md#EventServiceGetEventsForDay) | **Post** /event.EventService/GetEventsForDay | 
[**EventServiceGetEventsForMonth**](EventServiceApi.md#EventServiceGetEventsForMonth) | **Post** /event.EventService/GetEventsForMonth | 
[**EventServiceGetEventsForWeek**](EventServiceApi.md#EventServiceGetEventsForWeek) | **Post** /event.EventService/GetEventsForWeek | 
[**EventServiceUpdateEvent**](EventServiceApi.md#EventServiceUpdateEvent) | **Post** /event.EventService/UpdateEvent | 



## EventServiceCreateEvent

> EventCreateUpdateResponse EventServiceCreateEvent(ctx, body)



### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**body** | [**EventEvent**](EventEvent.md)|  | 

### Return type

[**EventCreateUpdateResponse**](eventCreateUpdateResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EventServiceDeleteEvent

> EventDeleteResponse EventServiceDeleteEvent(ctx, body)



### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**body** | [**EventDeleteRequest**](EventDeleteRequest.md)|  | 

### Return type

[**EventDeleteResponse**](eventDeleteResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EventServiceGetEventsForDay

> EventGetEventsResponse EventServiceGetEventsForDay(ctx, body)



### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**body** | [**EventGetEventsRequest**](EventGetEventsRequest.md)|  | 

### Return type

[**EventGetEventsResponse**](eventGetEventsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EventServiceGetEventsForMonth

> EventGetEventsResponse EventServiceGetEventsForMonth(ctx, body)



### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**body** | [**EventGetEventsRequest**](EventGetEventsRequest.md)|  | 

### Return type

[**EventGetEventsResponse**](eventGetEventsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EventServiceGetEventsForWeek

> EventGetEventsResponse EventServiceGetEventsForWeek(ctx, body)



### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**body** | [**EventGetEventsRequest**](EventGetEventsRequest.md)|  | 

### Return type

[**EventGetEventsResponse**](eventGetEventsResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EventServiceUpdateEvent

> EventCreateUpdateResponse EventServiceUpdateEvent(ctx, body)



### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**body** | [**EventEvent**](EventEvent.md)|  | 

### Return type

[**EventCreateUpdateResponse**](eventCreateUpdateResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

