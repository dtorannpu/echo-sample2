import * as z from "zod";

export type UUID = z.infer<ReturnType<typeof z.uuidv7>>;
