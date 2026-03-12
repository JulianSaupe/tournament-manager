export type ApiResult<T> = ApiSuccessResult<T> | ApiErrorResult;
export type ApiSuccessResult<T> = { success: true; data: T; error?: never };
export type ApiErrorResult = { success: false; data?: never; error: Error };

export const ApiSuccessResult = <T>(data: T): ApiSuccessResult<T> => ({ success: true, data });
export const ApiErrorResult = (error: Error): ApiErrorResult => ({ success: false, error });
