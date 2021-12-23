package com.example.demo.entity;

import com.fasterxml.jackson.annotation.JsonProperty;

public class ArcAddParam {
    private String title;
    private String cover;
    private int tid;
    @JsonProperty("no_reprint")
    private int noReprint;
    private String desc;
    private String tag;
    private int copyright;
    private String source;

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public String getCover() {
        return cover;
    }

    public void setCover(String cover) {
        this.cover = cover;
    }

    public int getTid() {
        return tid;
    }

    public void setTid(int tid) {
        this.tid = tid;
    }

    public int getNoReprint() {
        return noReprint;
    }

    public void setNoReprint(int noReprint) {
        this.noReprint = noReprint;
    }

    public String getDesc() {
        return desc;
    }

    public void setDesc(String desc) {
        this.desc = desc;
    }

    public String getTag() {
        return tag;
    }

    public void setTag(String tag) {
        this.tag = tag;
    }

    public int getCopyright() {
        return copyright;
    }

    public void setCopyright(int copyright) {
        this.copyright = copyright;
    }

    public String getSource() {
        return source;
    }

    public void setSource(String source) {
        this.source = source;
    }
}
