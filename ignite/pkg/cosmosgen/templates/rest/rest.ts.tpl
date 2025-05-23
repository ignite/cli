// FOR REFERENCE ONLY -- NOT A TEMPLATE

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, ResponseType } from "axios";
import { QueryAllBalancesResponse, QueryBalanceResponse, QueryDenomMetadataByQueryStringResponse, QueryDenomMetadataResponse, QueryDenomOwnersByQueryResponse, QueryDenomOwnersResponse, QueryDenomsMetadataResponse, QueryParamsResponse, QuerySpendableBalancesResponse, QuerySupplyOfResponse, QueryTotalSupplyResponse } from "./module";

export type QueryParamsType = Record<string | number, any>;

export interface FullRequestParams extends Omit<AxiosRequestConfig, "data" | "params" | "url" | "responseType"> {
  /** set parameter to `true` for call `securityWorker` for this request */
  secure?: boolean;
  /** request path */
  path: string;
  /** content type of request body */
  type?: ContentType;
  /** query params */
  query?: QueryParamsType;
  /** format of response (i.e. response.json() -> format: "json") */
  format?: ResponseType;
  /** request body */
  body?: unknown;
}

export type RequestParams = Omit<FullRequestParams, "body" | "method" | "query" | "path">;

export interface ApiConfig<SecurityDataType = unknown> extends Omit<AxiosRequestConfig, "data" | "cancelToken"> {
  securityWorker?: (
    securityData: SecurityDataType | null,
  ) => Promise<AxiosRequestConfig | void> | AxiosRequestConfig | void;
  secure?: boolean;
  format?: ResponseType;
}

export enum ContentType {
  Json = "application/json",
  FormData = "multipart/form-data",
  UrlEncoded = "application/x-www-form-urlencoded",
}

export class HttpClient<SecurityDataType = unknown> {
  public instance: AxiosInstance;
  private securityData: SecurityDataType | null = null;
  private securityWorker?: ApiConfig<SecurityDataType>["securityWorker"];
  private secure?: boolean;
  private format?: ResponseType;

  constructor({ securityWorker, secure, format, ...axiosConfig }: ApiConfig<SecurityDataType> = {}) {
    this.instance = axios.create({ ...axiosConfig, baseURL: axiosConfig.baseURL || "" });
    this.secure = secure;
    this.format = format;
    this.securityWorker = securityWorker;
  }

  public setSecurityData = (data: SecurityDataType | null) => {
    this.securityData = data;
  };

  private mergeRequestParams(params1: AxiosRequestConfig, params2?: AxiosRequestConfig): AxiosRequestConfig {
    return {
      ...this.instance.defaults,
      ...params1,
      ...(params2 || {}),
      headers: {
        ...(this.instance.defaults.headers ),
        ...(params1.headers || {}),
        ...((params2 && params2.headers) || {}),
      },
    } as AxiosRequestConfig;
  }

  private createFormData(input: Record<string, unknown>): FormData {
    return Object.keys(input || {}).reduce((formData, key) => {
      const property = input[key];
      formData.append(
        key,
        property instanceof Blob
          ? property
          : typeof property === "object" && property !== null
          ? JSON.stringify(property)
          : `${property}`,
      );
      return formData;
    }, new FormData());
  }

  public request = async <T = any>({
    secure,
    path,
    type,
    query,
    format,
    body,
    ...params
  }: FullRequestParams): Promise<AxiosResponse<T>> => {
    const secureParams =
      ((typeof secure === "boolean" ? secure : this.secure) &&
        this.securityWorker &&
        (await this.securityWorker(this.securityData))) ||
      {};
    const requestParams = this.mergeRequestParams(params, secureParams);
    const responseFormat = (format && this.format) || void 0;

    if (type === ContentType.FormData && body && body !== null && typeof body === "object") {
      requestParams.headers.common = { Accept: "*/*" };
      requestParams.headers.post = {};
      requestParams.headers.put = {};

      body = this.createFormData(body as Record<string, unknown>);
    }

    return this.instance.request({
      ...requestParams,
      headers: {
        ...(type && type !== ContentType.FormData ? { "Content-Type": type } : {}),
        ...(requestParams.headers || {}),
      },
      params: query,
      responseType: responseFormat,
      data: body,
      url: path,
    });
  };
}

/**
 * @title cosmos/bank/v1beta1/authz.proto
 * @version version not set
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * No description
   *
   * @tags Query
   * @name QueryAllBalances
   * @summary AllBalances queries the balance of all coins for a single account.
   * @request GET:/cosmos/bank/v1beta1/balances/{address}
   */
  queryAllBalances = (
    address: string,
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<QueryAllBalancesResponse>({
      path: `/cosmos/bank/v1beta1/balances/${address}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryBalance
   * @summary Balance queries the balance of a single coin for a single account.
   * @request GET:/cosmos/bank/v1beta1/balances/{address}/by_denom
   */
  queryBalance = (address: string, query?: { denom?: string }, params: RequestParams = {}) =>
    this.request<QueryBalanceResponse>({
      path: `/cosmos/bank/v1beta1/balances/${address}/by_denom`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
 * @description Since: cosmos-sdk 0.46
 * 
 * @tags Query
 * @name QueryDenomOwners
 * @summary DenomOwners queries for all account addresses that own a particular token
denomination.
 * @request GET:/cosmos/bank/v1beta1/denom_owners/{denom}
 */
  queryDenomOwners = (
    denom: string,
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<QueryDenomOwnersResponse>({
      path: `/cosmos/bank/v1beta1/denom_owners/${denom}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
 * @description Since: cosmos-sdk 0.50.3
 * 
 * @tags Query
 * @name QueryDenomOwnersByQuery
 * @summary DenomOwners queries for all account addresses that own a particular token
denomination.
 * @request GET:/cosmos/bank/v1beta1/denom_owners_by_query
 */
  queryDenomOwnersByQuery = (
    query?: {
      "denom": string;
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<QueryDenomOwnersByQueryResponse>({
      path: `/cosmos/bank/v1beta1/denom_owners_by_query`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
 * No description
 * 
 * @tags Query
 * @name QueryDenomsMetadata
 * @summary DenomsMetadata queries the client metadata for all registered coin
denominations.
 * @request GET:/cosmos/bank/v1beta1/denoms_metadata
 */
  queryDenomsMetadata = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<QueryDenomsMetadataResponse>({
      path: `/cosmos/bank/v1beta1/denoms_metadata`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryDenomMetadata
   * @summary DenomsMetadata queries the client metadata of a given coin denomination.
   * @request GET:/cosmos/bank/v1beta1/denoms_metadata/{denom}
   */
  queryDenomMetadata = (denom: string, params: RequestParams = {}) =>
    this.request<QueryDenomMetadataResponse>({
      path: `/cosmos/bank/v1beta1/denoms_metadata/${denom}`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryDenomMetadataByQueryString
   * @summary DenomsMetadata queries the client metadata of a given coin denomination.
   * @request GET:/cosmos/bank/v1beta1/denoms_metadata_by_query_string
   */
  queryDenomMetadataByQueryString = (
     query?: {
      "denom": string;
    },
    params: RequestParams = {}) =>
    this.request<QueryDenomMetadataByQueryStringResponse>({
      path: `/cosmos/bank/v1beta1/denoms_metadata_by_query_string`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryParams
   * @summary Params queries the parameters of x/bank module.
   * @request GET:/cosmos/bank/v1beta1/params
   */
  queryParams = (params: RequestParams = {}) =>
    this.request<QueryParamsResponse>({
      path: `/cosmos/bank/v1beta1/params`,
      method: "GET",
      format: "json",
      ...params,
    });

  /**
 * @description Since: cosmos-sdk 0.46
 * 
 * @tags Query
 * @name QuerySpendableBalances
 * @summary SpendableBalances queries the spenable balance of all coins for a single
account.
 * @request GET:/cosmos/bank/v1beta1/spendable_balances/{address}
 */
  querySpendableBalances = (
    address: string,
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<QuerySpendableBalancesResponse>({
      path: `/cosmos/bank/v1beta1/spendable_balances/${address}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QueryTotalSupply
   * @summary TotalSupply queries the total supply of all coins.
   * @request GET:/cosmos/bank/v1beta1/supply
   */
  queryTotalSupply = (
    query?: {
      "pagination.key"?: string;
      "pagination.offset"?: string;
      "pagination.limit"?: string;
      "pagination.count_total"?: boolean;
      "pagination.reverse"?: boolean;
    },
    params: RequestParams = {},
  ) =>
    this.request<QueryTotalSupplyResponse>({
      path: `/cosmos/bank/v1beta1/supply`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });

  /**
   * No description
   *
   * @tags Query
   * @name QuerySupplyOf
   * @summary SupplyOf queries the supply of a single coin.
   * @request GET:/cosmos/bank/v1beta1/supply/by_denom
   */
  querySupplyOf = (query?: { denom?: string }, params: RequestParams = {}) =>
    this.request<QuerySupplyOfResponse>({
      path: `/cosmos/bank/v1beta1/supply/by_denom`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
}