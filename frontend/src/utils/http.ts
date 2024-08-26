import { Ref, ref } from "vue";

// TODO: pagination
export const useFetch = async <T>(
  url: string,
  init?: RequestInit | undefined,
  isJson = true
) => {
  const data: Ref<T> = ref(new (Object as any)());
  const err: Ref<any> = ref(null);
  await fetch(url, init)
    .then((res) => {
      if (res.status != 200) err.value = res.statusText;
      else return isJson ? res.json() : res.text();
    })
    .then((json) => (data.value = json))
    .catch((err) => (err.value = err));
  return { data, err };
};

export const useHttp =
  (baseUrl: string) =>
  (token = "") => {
    const get = <T>(url: string, isJson = true) =>
      useFetch<T>(
        `${baseUrl}${url}`,
        {
          method: "GET",
          headers: isJson
            ? {
                "Content-Type": "application/json",
                Authorization: `${token}`,
              }
            : {},
        },
        isJson
      );
    const post = <T>(url: string, data: any) =>
      useFetch<T>(`${baseUrl}${url}`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `${token}`,
        },
        body: JSON.stringify(data),
      });
    const del = <T>(url: string, id: number) =>
      useFetch<T>(`${baseUrl}${url}/${id}`, {
        method: "DELETE",
        headers: {
          Authorization: `${token}`,
        },
      });
    const put = <T>(url: string, id: number, data: any) =>
      useFetch<T>(`${baseUrl}${url}/${id}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `${token}`,
        },
        body: JSON.stringify(data),
      });
    return { get, post, del, put };
  };

const baseUrl = "/api/v1";

export const http = useHttp(baseUrl)();

export type Pagination = {
  pageNum: number;
  pageSize: number;
  total: number;
};
export const getDataArr = <T>(
  api: string,
  _pagination = <Pagination>{ pageNum: 1, pageSize: 25, total: 25 }
) => {
  const data: Ref<T[]> = ref([]);
  const get = async (): Promise<T[]> => {
    const { data, err } = await http.get(api);
    if (err.value != null) return err.value;
    return (data as Ref<T[]>).value;
  };
  return { data, get };
};
