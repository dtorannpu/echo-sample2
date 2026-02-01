import * as z from "zod";

const schema = z.uuid();
export type UUID = z.infer<typeof schema>;
