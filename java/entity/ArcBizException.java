package com.example.demo.entity;

public class ArcBizException extends Exception {
    public ArcBizException(CommonReply cr) {
        super(String.format("business error: code(%s), message(%s)", cr.getCode(), cr.getMessage()));
    }
}
