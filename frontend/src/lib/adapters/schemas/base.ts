import {z} from "zod";

export const BaseApiSchema = <T extends z.ZodTypeAny>(dataSchema: T) =>
    z.object({
        data: dataSchema
    });