import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, ResponseType } from "axios";
import { QuerySimpleResponse } from "./types/planet/mars/mars";
import { QuerySimpleParamsResponse } from "./types/planet/mars/mars";
import { QueryWithPaginationResponse } from "./types/planet/mars/mars";
import { QueryWithQueryParamsResponse } from "./types/planet/mars/mars";
import { QueryWithQueryParamsWithPaginationResponse } from "./types/planet/mars/mars";

import type {SnakeCasedPropertiesDeep} from 'type-fest';

export type QueryParamsType = Record<string | number, any>;


export type ChangeProtoToJSPrimitives<T extends object> = {
  [key in keyof T]: T[key] extends Uint8Array | Date ? string :  T[key] extends object ? ChangeProtoToJSPrimitives<T[key]>: T[key];
  // ^^^^ This line is used to convert Uint8Array to string, if you want to keep Uint8Array as is, you can remove this line
}

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
 * @title ignite.planet.mars
 */
export class Api<SecurityDataType extends unknown> extends HttpClient<SecurityDataType> {
  /**
   * QueryQuerySimple
   *
   * @tags Query
   * @name queryQuerySimple
   * @request GET:/ignite/mars/query_simple
   */
  queryQuerySimple = (
    query?: Record<string, any>,
    params: RequestParams = {},
  ) =>
    this.request<SnakeCasedPropertiesDeep<ChangeProtoToJSPrimitives<QuerySimpleResponse>>>({
      path: `/ignite/mars/query_simple`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
  
  /**
   * QueryQuerySimpleParams
   *
   * @tags Query
   * @name queryQuerySimpleParams
   * @request GET:/ignite/mars/query_simple/{mytypefield}
   */
  queryQuerySimpleParams = (mytypefield: string,
    query?: Record<string, any>,
    params: RequestParams = {},
  ) =>
    this.request<SnakeCasedPropertiesDeep<ChangeProtoToJSPrimitives<QuerySimpleParamsResponse>>>({
      path: `/ignite/mars/query_simple/${mytypefield}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
  
  /**
   * QueryQueryParamsWithPagination
   *
   * @tags Query
   * @name queryQueryParamsWithPagination
   * @request GET:/ignite/mars/query_with_params/{mytypefield}
   */
  queryQueryParamsWithPagination = (mytypefield: string,
    query?: {
      "pagination"?: any /* TODO */;
    },
    params: RequestParams = {},
  ) =>
    this.request<SnakeCasedPropertiesDeep<ChangeProtoToJSPrimitives<QueryWithPaginationResponse>>>({
      path: `/ignite/mars/query_with_params/${mytypefield}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
  
  /**
   * QueryQueryWithQueryParams
   *
   * @tags Query
   * @name queryQueryWithQueryParams
   * @request GET:/ignite/mars/query_with_query_params/{mytypefield}
   */
  queryQueryWithQueryParams = (mytypefield: string,
    query?: {
      "mybool"?: boolean;
      "myrepeatedbool"?: boolean[];
      "query_param"?: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<SnakeCasedPropertiesDeep<ChangeProtoToJSPrimitives<QueryWithQueryParamsResponse>>>({
      path: `/ignite/mars/query_with_query_params/${mytypefield}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
  
  /**
   * QueryQueryWithQueryParamsWithPagination
   *
   * @tags Query
   * @name queryQueryWithQueryParamsWithPagination
   * @request GET:/ignite/mars/query_with_query_params/{mytypefield}
   */
  queryQueryWithQueryParamsWithPagination = (mytypefield: string,
    query?: {
      "pagination"?: any /* TODO */;
      "query_param"?: string;
    },
    params: RequestParams = {},
  ) =>
    this.request<SnakeCasedPropertiesDeep<ChangeProtoToJSPrimitives<QueryWithQueryParamsWithPaginationResponse>>>({
      path: `/ignite/mars/query_with_query_params/${mytypefield}`,
      method: "GET",
      query: query,
      format: "json",
      ...params,
    });
  
}