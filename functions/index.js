const {Datastore} = require('@google-cloud/datastore');

const datastore = new Datastore();

const ancestorKind = 'Environment';
const kind = 'Microservice';

/**
 * Responds to any HTTP request.
 *
 * @param {!express:Request} req HTTP request context.
 * @param {!express:Response} res HTTP response context.
 */
exports.appInventory = (req, res) => {
  let year = req.query.year;
  console.log(year);
  
  listMicroservices(year).then(result => {
    res.status(200).send(result);
  });
};

async function listMicroservices(year) {
  const ancestorKey = datastore.key([ancestorKind, 'Cloud Functions']);
  let query = datastore.createQuery(kind).hasAncestor(ancestorKey);
  if(!isNaN(parseInt(year))) {
    query = query.filter('year', '=', parseInt(year, 10))
  }

  const [microservices] = await datastore.runQuery(query);

  let result = "<h1>Cloud Functions Applications Inventory</h1><ul>";  
  for (const microservice of microservices) {
    result += "<li>" + microservice['name'] + " (" + microservice['year'] + ")</li>"
  }
  result += "</ul>"
  
  return result;
}
