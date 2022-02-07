export default class Page {

    public open(path: string): Promise<string> {
        return browser.url(`/?#/${path}`)
    }
}
