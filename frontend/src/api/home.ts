
export type Feed = {
    version: string;
    author: Author[];
}
export type Author = {
    name: string;
    email: string;
    uri: string;
    description: string;
    article: Article[];
}
export type Article = {
    title: string;
    time: string;
    content: string;
    url: string;
}
