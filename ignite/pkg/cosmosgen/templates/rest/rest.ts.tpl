/* eslint-disable */
import axios from 'axios'
import * as qs from 'qs'

// this is used to derive the proper return types for query endpoints
export type BaseQueryClient = {
  queryBalances: (address: string, params?: any) => Promise<any>
}

export class Api {
  private axios: any
  private baseURL: string

  constructor({ baseURL }: { baseURL: string }) {
    this.baseURL = baseURL
    this.axios = axios.create({
      baseURL,
      timeout: 30000,
      paramsSerializer: function(params: any) {
        return qs.stringify(params, { arrayFormat: 'repeat' })
      }
    })
  }
  
  // common helper for most simple operations
  private async handleRequest(url: string, params?: any): Promise<any> {
    try {
      const response = await this.axios.get(url, { params })
      return response
    } catch (e: any) {
      if (e.response?.data) {
        console.error('Error in API request:', e.response.data)
      }
      throw e
    }
  }
  
  // Return URL for specific module endpoints
  public getModuleEndpoint(endpoint: string): string {
    return `${this.baseURL}/${endpoint}`
  }
  
  // Methods for specific module endpoints can be added here
  // The actual methods will be auto-generated from OpenAPI specs in a real implementation
}
