from tyk.decorators import *
from gateway import TykGateway as tyk


@Hook
def MyAuthMiddleware(request, session, metadata, spec):
    auth = request.get_header('Authorization')
    if not auth:
        auth = request.object.params.get('authorization', None)

    if auth == '47a0c79c427728b3df4af62b9228c8ae':
        session.rate = 1000.0
        session.per = 1.0
        metadata["token"] = auth
    return request, session, metadata


@Hook
def MyPostMiddleware(request, session, spec):
    tyk.log("This is my post middleware", "info")
    request.object.set_headers["x-tyk-request"] = "something"
    return request, session
