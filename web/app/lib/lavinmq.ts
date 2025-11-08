/**
 * Make a request to the LavinMQ API
 * @param url - The URL to make the request to
 * @param options - The options for the request
 * @returns The response from the request
 */
export const request = async <T>(
  url: string,
  options?: RequestInit,
  responseType: "json" | "text" = "json"
): Promise<T> => {
  const response = await fetch(process.env.LAVINMQ_API + url, {
    ...(options || {}),
    headers: {
      ...(options?.headers || {}),
      Authorization: `Basic ${btoa(process.env.LAVINMQ_USERNAME + ":" + process.env.LAVINMQ_PASSWORD)}`,
    },
  });

  if (responseType === "json") {
    return response.json() as T;
  } else if (responseType === "text") {
    return response.text() as T;
  }

  throw new Error("Invalid response type");
};
