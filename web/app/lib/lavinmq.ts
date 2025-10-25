/**
 * Make a request to the LavinMQ API
 * @param url - The URL to make the request to
 * @param options - The options for the request
 * @returns The response from the request
 */
export const request = async <T>(
  url: string,
  options?: RequestInit
): Promise<T> => {
  const response = await fetch(process.env.LAVINMQ_API + url, {
    ...(options || {}),
    headers: {
      ...(options?.headers || {}),
      Authorization: `Basic ${btoa(process.env.LAVINMQ_USERNAME + ":" + process.env.LAVINMQ_PASSWORD)}`,
    },
  });
  return response.json() as T;
};
